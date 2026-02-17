package gene

import (
	"fmt"
	"runtime"
	"time"
)

// SolidifyRequest contains the data needed to solidify an experience.
type SolidifyRequest struct {
	Signals    []string       // What triggered this event
	Strategy   []string       // What steps were taken
	Summary    string         // Human-readable description
	Outcome    CapsuleOutcome // Result of the action
	SensorData map[string]any // Sensor readings at time of event
	GeneID     string         // If an existing Gene was used (empty = new experience)
}

// SolidifyResult captures what happened during solidification.
type SolidifyResult struct {
	CapsuleCreated bool
	CapsuleID      string
	GeneUpdated    bool
	GeneCreated    bool
	GeneID         string
	EventID        string
}

// Solidifier packages learned experiences into Capsules and evolves Genes.
// Ported from Evolver's solidify.js.
type Solidifier struct {
	store           *Store
	signalExtractor *SignalExtractor
}

// NewSolidifier creates a new Solidifier.
func NewSolidifier(store *Store, se *SignalExtractor) *Solidifier {
	return &Solidifier{
		store:           store,
		signalExtractor: se,
	}
}

// Solidify processes a request and creates/updates Genes and Capsules.
// This is the core evolution mechanism.
func (s *Solidifier) Solidify(req SolidifyRequest) (*SolidifyResult, error) {
	now := time.Now()
	result := &SolidifyResult{}

	envFP := s.buildEnvFingerprint()

	// Create capsule
	capsuleID := fmt.Sprintf("cap_%d", now.UnixMilli())
	capsule := Capsule{
		ID:             capsuleID,
		GeneID:         req.GeneID,
		Trigger:        req.Signals,
		Summary:        req.Summary,
		Confidence:     s.computeCapsuleConfidence(req),
		Outcome:        req.Outcome,
		EnvFingerprint: envFP,
		SensorData:     req.SensorData,
		CreatedAt:      now.Format(time.RFC3339),
	}

	if err := s.store.AddCapsule(capsule); err != nil {
		return nil, fmt.Errorf("add capsule: %w", err)
	}
	result.CapsuleCreated = true
	result.CapsuleID = capsuleID

	// If an existing Gene was used, update its confidence
	if req.GeneID != "" {
		if err := s.updateGeneFromCapsule(req.GeneID, capsule); err != nil {
			return result, fmt.Errorf("update gene: %w", err)
		}
		result.GeneUpdated = true
		result.GeneID = req.GeneID
	} else if req.Outcome.Status == "success" && len(req.Strategy) > 0 {
		// Create a new Gene from successful novel handling
		newGene, err := s.createGeneFromExperience(req, envFP, now)
		if err != nil {
			return result, fmt.Errorf("create gene: %w", err)
		}
		result.GeneCreated = true
		result.GeneID = newGene.ID

		// Update capsule's gene reference
		capsule.GeneID = newGene.ID
	}

	// Record evolution event
	evtID := fmt.Sprintf("evt_%d", now.UnixMilli())
	intent := s.inferIntent(req)
	evt := EvolutionEvent{
		ID:        evtID,
		Intent:    intent,
		GeneID:    result.GeneID,
		CapsuleID: capsuleID,
		Signals:   req.Signals,
		Outcome:   req.Outcome,
		Timestamp: now.Format(time.RFC3339),
	}
	if err := s.store.AppendEvent(evt); err != nil {
		return result, fmt.Errorf("append event: %w", err)
	}
	result.EventID = evtID

	return result, nil
}

// updateGeneFromCapsule adjusts a Gene's confidence based on capsule outcome.
// Ported from Evolver's epigenetic marks concept.
func (s *Solidifier) updateGeneFromCapsule(geneID string, capsule Capsule) error {
	gene := s.store.GetGene(geneID)
	if gene == nil {
		return fmt.Errorf("gene %s not found", geneID)
	}

	// Adjust confidence based on outcome
	switch capsule.Outcome.Status {
	case "success":
		// Increase confidence, with diminishing returns
		boost := (1.0 - gene.Confidence) * 0.1
		gene.Confidence += boost
	case "failed":
		// Decrease confidence
		penalty := gene.Confidence * 0.15
		gene.Confidence -= penalty
	case "partial":
		// Slight increase
		boost := (1.0 - gene.Confidence) * 0.03
		gene.Confidence += boost
	}

	// Clamp confidence
	if gene.Confidence > 1.0 {
		gene.Confidence = 1.0
	}
	if gene.Confidence < 0.0 {
		gene.Confidence = 0.0
	}

	// Increment verified_by for successful capsules
	if capsule.Outcome.Status == "success" {
		gene.VerifiedBy++
	}

	return s.store.UpdateGene(*gene)
}

// createGeneFromExperience creates a new Gene from a successful novel experience.
func (s *Solidifier) createGeneFromExperience(req SolidifyRequest, envFP EnvFingerprint, now time.Time) (*Gene, error) {
	geneID := fmt.Sprintf("gene_learned_%d", now.UnixMilli())

	category := s.inferIntent(req)
	scenario := s.signalExtractor.ExtractScenarioHints()

	gene := Gene{
		ID:             geneID,
		Category:       category,
		Scenario:       scenario,
		SignalsMatch:   req.Signals,
		Strategy:       req.Strategy,
		Constraints:    GeneConstraints{MaxRetries: 3, TimeoutSeconds: 300},
		Confidence:     0.5, // Start at moderate confidence
		VerifiedBy:     1,   // Verified by this node
		EnvFingerprint: envFP,
		CreatedAt:      now.Format(time.RFC3339),
	}

	if err := s.store.AddGene(gene); err != nil {
		return nil, err
	}

	return &gene, nil
}

// computeCapsuleConfidence determines initial capsule confidence.
func (s *Solidifier) computeCapsuleConfidence(req SolidifyRequest) float64 {
	base := 0.5

	switch req.Outcome.Status {
	case "success":
		base = 0.8
	case "partial":
		base = 0.5
	case "failed":
		base = 0.2
	}

	// Having sensor data is a positive signal
	if len(req.SensorData) > 0 {
		base += 0.05
	}

	// Having a clear strategy is a positive signal
	if len(req.Strategy) >= 3 {
		base += 0.05
	}

	if base > 1.0 {
		base = 1.0
	}
	return base
}

// inferIntent determines whether this is a repair, optimize, or innovate action.
func (s *Solidifier) inferIntent(req SolidifyRequest) string {
	for _, sig := range req.Signals {
		switch sig {
		case SignalSensorError, SignalThresholdBreach:
			return "repair"
		case SignalResponseTooSlow, SignalStrategyProven, SignalTimePattern:
			return "optimize"
		case SignalNewPatternDetected, SignalUnknownSituation, SignalCrossSensorAnomaly:
			return "innovate"
		}
	}
	return "optimize"
}

// buildEnvFingerprint captures the current runtime environment.
func (s *Solidifier) buildEnvFingerprint() EnvFingerprint {
	return EnvFingerprint{
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}
