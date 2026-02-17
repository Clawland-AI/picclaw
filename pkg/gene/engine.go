package gene

import (
	"fmt"
	"strings"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// Engine is the top-level orchestrator for the Gene Evolution system.
// It ties together the Store, SignalExtractor, Selector, and Solidifier.
type Engine struct {
	store      *Store
	extractor  *SignalExtractor
	solidifier *Solidifier
	config     GeneConfig
}

// NewEngine creates and initializes a Gene Engine.
func NewEngine(workspace string, config GeneConfig) (*Engine, error) {
	store := NewStore(workspace)
	if err := store.Init(); err != nil {
		return nil, fmt.Errorf("init gene store: %w", err)
	}

	extractor := NewSignalExtractor(workspace)
	solidifier := NewSolidifier(store, extractor)

	engine := &Engine{
		store:      store,
		extractor:  extractor,
		solidifier: solidifier,
		config:     config,
	}

	logger.InfoCF("gene", "Gene engine initialized",
		map[string]interface{}{
			"genes":    store.GeneCount(),
			"capsules": store.CapsuleCount(),
			"strategy": config.Strategy,
		})

	return engine, nil
}

// SeedIfEmpty seeds the store with default genes if it has none.
func (e *Engine) SeedIfEmpty(seeds []Gene) error {
	return e.store.SeedGenes(seeds)
}

// GetGeneContext builds a formatted string of the top-matched Gene strategies
// for inclusion in the agent's system prompt.
// Called by ContextBuilder.BuildSystemPrompt().
func (e *Engine) GetGeneContext() string {
	signals := e.extractor.ExtractSignals()
	if len(signals) == 0 {
		return e.buildSummaryOnly()
	}

	preset := e.getPreset()
	genes := e.store.GetGenes()

	// Select top 3 matching genes
	scored := SelectGenes(genes, signals, preset, 3)
	if len(scored) == 0 {
		return e.buildSummaryWithSignals(signals)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Active signals: %s\n\n", strings.Join(signals, ", ")))

	for i, sg := range scored {
		sb.WriteString(fmt.Sprintf("### Strategy %d: %s (confidence: %.0f%%, score: %.2f)\n",
			i+1, sg.Gene.ID, sg.Gene.Confidence*100, sg.Score))
		sb.WriteString(fmt.Sprintf("Category: %s | Scenario: %s\n", sg.Gene.Category, sg.Gene.Scenario))
		sb.WriteString("Steps:\n")
		for j, step := range sg.Gene.Strategy {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", j+1, step))
		}

		// Include relevant capsule if available
		capsules := e.store.GetCapsules()
		if cap := SelectCapsule(capsules, sg.Gene.ID, signals); cap != nil {
			sb.WriteString(fmt.Sprintf("\nPrevious experience: %s (outcome: %s)\n", cap.Summary, cap.Outcome.Status))
		}
		sb.WriteString("\n")
	}

	// Add drift info for transparency
	drift := ComputeDriftIntensity(e.store.GeneCount(), e.store.CapsuleCount())
	if drift > 0.7 {
		sb.WriteString("⚡ High drift intensity detected — consider innovative approaches for novel patterns.\n")
	}

	return sb.String()
}

// ConsiderSolidify evaluates whether to solidify the result of an action.
// Called by the agent loop after successful handling.
func (e *Engine) ConsiderSolidify(signals []string, strategy []string, summary string, outcomeStatus string, sensorData map[string]any) *SolidifyResult {
	if len(signals) == 0 && len(strategy) == 0 {
		return nil
	}

	// Check if an existing gene was matched
	preset := e.getPreset()
	genes := e.store.GetGenes()
	scored := SelectGenes(genes, signals, preset, 1)

	var geneID string
	if len(scored) > 0 {
		geneID = scored[0].Gene.ID
	}

	req := SolidifyRequest{
		Signals:    signals,
		Strategy:   strategy,
		Summary:    summary,
		Outcome:    CapsuleOutcome{Status: outcomeStatus, Details: summary},
		SensorData: sensorData,
		GeneID:     geneID,
	}

	result, err := e.solidifier.Solidify(req)
	if err != nil {
		logger.ErrorCF("gene", "Solidify failed",
			map[string]interface{}{
				"error":   err.Error(),
				"signals": signals,
			})
		return nil
	}

	logger.InfoCF("gene", "Experience solidified",
		map[string]interface{}{
			"capsule_id":   result.CapsuleID,
			"gene_id":      result.GeneID,
			"gene_created": result.GeneCreated,
			"gene_updated": result.GeneUpdated,
		})

	return result
}

// ReportGene is called when the agent explicitly reports a learned experience
// via the report_gene tool.
func (e *Engine) ReportGene(summary string, strategy []string, signals []string, outcome string) (*SolidifyResult, error) {
	req := SolidifyRequest{
		Signals:  signals,
		Strategy: strategy,
		Summary:  summary,
		Outcome:  CapsuleOutcome{Status: outcome, Details: summary},
	}

	return e.solidifier.Solidify(req)
}

// GetStats returns current gene engine statistics.
func (e *Engine) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"genes":    e.store.GeneCount(),
		"capsules": e.store.CapsuleCount(),
		"strategy": e.config.Strategy,
		"drift":    ComputeDriftIntensity(e.store.GeneCount(), e.store.CapsuleCount()),
	}
}

// ExtractSignalsFromText exposes the signal extractor for external use.
func (e *Engine) ExtractSignalsFromText(text string) []string {
	return e.extractor.ExtractSignalsFromText(text)
}

// getPreset returns the active strategy preset.
func (e *Engine) getPreset() StrategyPreset {
	if preset, ok := StrategyPresets[e.config.Strategy]; ok {
		// Apply overrides from config
		if e.config.MinConfidence > 0 {
			preset.MinConfidence = e.config.MinConfidence
		}
		if e.config.MinVerifiedBy > 0 {
			preset.MinVerifiedBy = e.config.MinVerifiedBy
		}
		return preset
	}
	return StrategyPresets["balanced"]
}

// buildSummaryOnly returns a brief gene status when no signals are active.
func (e *Engine) buildSummaryOnly() string {
	gc := e.store.GeneCount()
	cc := e.store.CapsuleCount()
	if gc == 0 {
		return ""
	}
	return fmt.Sprintf("Gene pool: %d genes, %d capsules. Strategy: %s. No active signals.\n", gc, cc, e.config.Strategy)
}

// buildSummaryWithSignals returns status with signals but no matching genes.
func (e *Engine) buildSummaryWithSignals(signals []string) string {
	return fmt.Sprintf("Active signals: %s\nNo matching genes found. Consider creating a new strategy using the report_gene tool after handling this situation.\n",
		strings.Join(signals, ", "))
}
