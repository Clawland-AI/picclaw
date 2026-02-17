package gene

import (
	"math"
	"sort"
	"strings"
)

// ScoredGene pairs a Gene with its match score against current signals.
type ScoredGene struct {
	Gene  Gene
	Score float64
}

// ScoreGene computes a match score between a Gene's signals_match and the
// current signals. Ported from Evolver's selector.js scoreGene().
//
// Score factors:
//   - Signal match ratio (how many of gene's signals are present)
//   - Confidence weight (higher confidence = higher score)
//   - Verification bonus (more verifications = more trustworthy)
func ScoreGene(g Gene, signals []string) float64 {
	if len(g.SignalsMatch) == 0 {
		return 0
	}

	// Compute signal match ratio
	signalSet := make(map[string]struct{}, len(signals))
	for _, s := range signals {
		signalSet[strings.ToLower(s)] = struct{}{}
	}

	matched := 0
	for _, required := range g.SignalsMatch {
		if _, ok := signalSet[strings.ToLower(required)]; ok {
			matched++
		}
	}

	if matched == 0 {
		return 0
	}

	matchRatio := float64(matched) / float64(len(g.SignalsMatch))

	// Confidence contributes to score
	confidenceWeight := g.Confidence * 0.3

	// Verification count as log bonus (diminishing returns)
	verificationBonus := 0.0
	if g.VerifiedBy > 0 {
		verificationBonus = math.Log2(float64(g.VerifiedBy)+1) * 0.1
	}

	// Combined score: match ratio is primary (70%), confidence (20%), verification (10%)
	score := matchRatio*0.7 + confidenceWeight + verificationBonus

	return math.Min(score, 1.0)
}

// SelectGenes selects the top-N genes that match the given signals,
// filtered by the strategy preset. Ported from Evolver's selectGene().
func SelectGenes(genes []Gene, signals []string, preset StrategyPreset, maxResults int) []ScoredGene {
	var candidates []ScoredGene

	for _, g := range genes {
		// Apply strategy filters
		if !preset.AllowUnverified && g.VerifiedBy < preset.MinVerifiedBy {
			continue
		}
		if g.Confidence < preset.MinConfidence {
			continue
		}
		if preset.CategoryFilter != "" && g.Category != preset.CategoryFilter {
			continue
		}

		score := ScoreGene(g, signals)
		if score > 0 {
			candidates = append(candidates, ScoredGene{Gene: g, Score: score})
		}
	}

	// Sort by score descending
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	if maxResults > 0 && len(candidates) > maxResults {
		candidates = candidates[:maxResults]
	}

	return candidates
}

// ComputeDriftIntensity calculates how much genetic drift is present.
// Small populations (few genes on a single node) = higher drift.
// Ported from Evolver's computeDriftIntensity().
//
// Returns a value between 0.0 (no drift) and 1.0 (maximum drift).
// High drift means the node should be more exploratory in creating new genes.
func ComputeDriftIntensity(geneCount int, capsuleCount int) float64 {
	if geneCount == 0 {
		return 1.0 // No genes at all = maximum drift, need innovation
	}

	// Small gene pool = higher drift
	geneFactor := 1.0 / (1.0 + math.Log2(float64(geneCount)))

	// Low capsule-to-gene ratio = less experience, higher drift
	ratio := float64(capsuleCount) / float64(geneCount)
	experienceFactor := 1.0 / (1.0 + ratio*0.5)

	drift := (geneFactor + experienceFactor) / 2.0
	return math.Min(drift, 1.0)
}

// SelectCapsule finds the most relevant capsule for a given gene and signals.
// Prioritizes recent capsules with matching triggers.
// Ported from Evolver's selectCapsule().
func SelectCapsule(capsules []Capsule, geneID string, signals []string) *Capsule {
	signalSet := make(map[string]struct{}, len(signals))
	for _, s := range signals {
		signalSet[strings.ToLower(s)] = struct{}{}
	}

	var bestCapsule *Capsule
	bestScore := 0.0

	for i, c := range capsules {
		if c.GeneID != geneID {
			continue
		}

		// Score by trigger match + outcome
		matched := 0
		for _, t := range c.Trigger {
			if _, ok := signalSet[strings.ToLower(t)]; ok {
				matched++
			}
		}

		if matched == 0 {
			continue
		}

		score := float64(matched) / float64(len(c.Trigger))

		// Prefer successful capsules
		if c.Outcome.Status == "success" {
			score *= 1.2
		}

		// Prefer higher confidence
		score += c.Confidence * 0.1

		if score > bestScore {
			bestScore = score
			bestCapsule = &capsules[i]
		}
	}

	return bestCapsule
}
