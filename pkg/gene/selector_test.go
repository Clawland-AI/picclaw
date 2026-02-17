package gene

import (
	"testing"
)

func makeTestGene(id string, signals []string, confidence float64, verifiedBy int, category string) Gene {
	return Gene{
		ID:           id,
		Category:     category,
		Scenario:     "datacenter",
		SignalsMatch: signals,
		Strategy:     []string{"step1"},
		Confidence:   confidence,
		VerifiedBy:   verifiedBy,
	}
}

func TestScoreGene_NoSignalsMatch_ReturnsZero(t *testing.T) {
	g := makeTestGene("g1", []string{"sensor_error"}, 0.8, 3, "repair")
	score := ScoreGene(g, []string{"threshold_breach"})
	if score != 0 {
		t.Errorf("expected 0 for no matching signals, got %f", score)
	}
}

func TestScoreGene_EmptyGeneSignals_ReturnsZero(t *testing.T) {
	g := makeTestGene("g1", []string{}, 0.8, 3, "repair")
	score := ScoreGene(g, []string{"sensor_error"})
	if score != 0 {
		t.Errorf("expected 0 for empty gene signals, got %f", score)
	}
}

func TestScoreGene_FullMatch_HighScore(t *testing.T) {
	g := makeTestGene("g1", []string{"sensor_error", "threshold_breach"}, 0.9, 5, "repair")
	score := ScoreGene(g, []string{"sensor_error", "threshold_breach"})

	if score <= 0 {
		t.Errorf("expected positive score for full match, got %f", score)
	}
	if score > 1.0 {
		t.Errorf("score should be capped at 1.0, got %f", score)
	}
}

func TestScoreGene_PartialMatch_LowerScore(t *testing.T) {
	g := makeTestGene("g1", []string{"sensor_error", "threshold_breach", "time_pattern"}, 0.8, 3, "repair")
	scoreFull := ScoreGene(g, []string{"sensor_error", "threshold_breach", "time_pattern"})
	scorePartial := ScoreGene(g, []string{"sensor_error"})

	if scorePartial >= scoreFull {
		t.Errorf("partial match score (%f) should be less than full match (%f)", scorePartial, scoreFull)
	}
}

func TestScoreGene_CaseInsensitive(t *testing.T) {
	g := makeTestGene("g1", []string{"Sensor_Error"}, 0.8, 3, "repair")
	score := ScoreGene(g, []string{"sensor_error"})

	if score == 0 {
		t.Error("signal matching should be case-insensitive")
	}
}

func TestScoreGene_HighConfidence_HigherScore(t *testing.T) {
	gHigh := makeTestGene("g1", []string{"sensor_error"}, 0.95, 0, "repair")
	gLow := makeTestGene("g2", []string{"sensor_error"}, 0.1, 0, "repair")

	scoreHigh := ScoreGene(gHigh, []string{"sensor_error"})
	scoreLow := ScoreGene(gLow, []string{"sensor_error"})

	if scoreHigh <= scoreLow {
		t.Errorf("higher confidence should yield higher score: high=%f, low=%f", scoreHigh, scoreLow)
	}
}

func TestScoreGene_VerificationBonus(t *testing.T) {
	gVerified := makeTestGene("g1", []string{"sensor_error"}, 0.8, 10, "repair")
	gUnverified := makeTestGene("g2", []string{"sensor_error"}, 0.8, 0, "repair")

	scoreV := ScoreGene(gVerified, []string{"sensor_error"})
	scoreU := ScoreGene(gUnverified, []string{"sensor_error"})

	if scoreV <= scoreU {
		t.Errorf("verified gene should score higher: verified=%f, unverified=%f", scoreV, scoreU)
	}
}

func TestSelectGenes_Filtering(t *testing.T) {
	genes := []Gene{
		makeTestGene("g1", []string{"sensor_error"}, 0.9, 5, "repair"),
		makeTestGene("g2", []string{"sensor_error"}, 0.3, 0, "repair"),   // low confidence
		makeTestGene("g3", []string{"sensor_error"}, 0.9, 1, "optimize"), // wrong category for filter
	}

	preset := StrategyPreset{
		MinConfidence:   0.7,
		MinVerifiedBy:   2,
		AllowUnverified: false,
		CategoryFilter:  "",
	}

	results := SelectGenes(genes, []string{"sensor_error"}, preset, 10)

	if len(results) != 1 {
		t.Errorf("expected 1 gene after filtering, got %d", len(results))
	}
	if len(results) > 0 && results[0].Gene.ID != "g1" {
		t.Errorf("expected g1 to be selected, got %s", results[0].Gene.ID)
	}
}

func TestSelectGenes_CategoryFilter(t *testing.T) {
	genes := []Gene{
		makeTestGene("g1", []string{"sensor_error"}, 0.9, 5, "repair"),
		makeTestGene("g2", []string{"sensor_error"}, 0.9, 5, "optimize"),
	}

	preset := StrategyPreset{
		MinConfidence:   0.5,
		MinVerifiedBy:   1,
		AllowUnverified: false,
		CategoryFilter:  "repair",
	}

	results := SelectGenes(genes, []string{"sensor_error"}, preset, 10)

	if len(results) != 1 {
		t.Errorf("expected 1 gene with category filter, got %d", len(results))
	}
	if len(results) > 0 && results[0].Gene.ID != "g1" {
		t.Errorf("expected g1, got %s", results[0].Gene.ID)
	}
}

func TestSelectGenes_MaxResults(t *testing.T) {
	genes := []Gene{
		makeTestGene("g1", []string{"sensor_error"}, 0.9, 5, "repair"),
		makeTestGene("g2", []string{"sensor_error"}, 0.85, 5, "repair"),
		makeTestGene("g3", []string{"sensor_error"}, 0.8, 5, "repair"),
	}

	preset := StrategyPreset{
		MinConfidence:   0.5,
		MinVerifiedBy:   1,
		AllowUnverified: false,
	}

	results := SelectGenes(genes, []string{"sensor_error"}, preset, 2)

	if len(results) != 2 {
		t.Errorf("expected 2 results with maxResults=2, got %d", len(results))
	}
}

func TestSelectGenes_SortedByScore(t *testing.T) {
	genes := []Gene{
		makeTestGene("low", []string{"sensor_error"}, 0.7, 2, "repair"),
		makeTestGene("high", []string{"sensor_error"}, 0.95, 10, "repair"),
	}

	preset := StrategyPreset{
		MinConfidence:   0.5,
		MinVerifiedBy:   1,
		AllowUnverified: false,
	}

	results := SelectGenes(genes, []string{"sensor_error"}, preset, 10)

	if len(results) < 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Gene.ID != "high" {
		t.Errorf("expected highest-scoring gene first, got %s", results[0].Gene.ID)
	}
}

func TestSelectGenes_AllowUnverified(t *testing.T) {
	genes := []Gene{
		makeTestGene("g1", []string{"sensor_error"}, 0.5, 0, "repair"),
	}

	strict := StrategyPreset{
		MinConfidence:   0.3,
		MinVerifiedBy:   1,
		AllowUnverified: false,
	}
	lenient := StrategyPreset{
		MinConfidence:   0.3,
		MinVerifiedBy:   0,
		AllowUnverified: true,
	}

	resultsStrict := SelectGenes(genes, []string{"sensor_error"}, strict, 10)
	resultsLenient := SelectGenes(genes, []string{"sensor_error"}, lenient, 10)

	if len(resultsStrict) != 0 {
		t.Errorf("strict preset should filter out unverified gene, got %d", len(resultsStrict))
	}
	if len(resultsLenient) != 1 {
		t.Errorf("lenient preset should include unverified gene, got %d", len(resultsLenient))
	}
}

func TestComputeDriftIntensity_NoGenes(t *testing.T) {
	drift := ComputeDriftIntensity(0, 0)
	if drift != 1.0 {
		t.Errorf("expected 1.0 for no genes, got %f", drift)
	}
}

func TestComputeDriftIntensity_ManyGenes(t *testing.T) {
	drift := ComputeDriftIntensity(100, 500)
	if drift >= 1.0 {
		t.Errorf("expected drift < 1.0 for large gene pool, got %f", drift)
	}
	if drift <= 0 {
		t.Errorf("drift should be positive, got %f", drift)
	}
}

func TestComputeDriftIntensity_FewGenesHigherThanMany(t *testing.T) {
	driftFew := ComputeDriftIntensity(2, 1)
	driftMany := ComputeDriftIntensity(50, 200)

	if driftFew <= driftMany {
		t.Errorf("fewer genes should produce higher drift: few=%f, many=%f", driftFew, driftMany)
	}
}

func TestSelectCapsule_MatchesGeneAndSignals(t *testing.T) {
	capsules := []Capsule{
		{GeneID: "gene_a", Trigger: []string{"sensor_error"}, Outcome: CapsuleOutcome{Status: "success"}, Confidence: 0.8},
		{GeneID: "gene_b", Trigger: []string{"threshold_breach"}, Outcome: CapsuleOutcome{Status: "success"}, Confidence: 0.9},
	}

	result := SelectCapsule(capsules, "gene_a", []string{"sensor_error"})
	if result == nil {
		t.Fatal("expected a capsule match")
	}
	if result.GeneID != "gene_a" {
		t.Errorf("expected gene_a capsule, got %s", result.GeneID)
	}
}

func TestSelectCapsule_NoMatch(t *testing.T) {
	capsules := []Capsule{
		{GeneID: "gene_a", Trigger: []string{"sensor_error"}, Outcome: CapsuleOutcome{Status: "success"}},
	}

	result := SelectCapsule(capsules, "gene_b", []string{"sensor_error"})
	if result != nil {
		t.Error("expected nil for non-matching gene ID")
	}

	result = SelectCapsule(capsules, "gene_a", []string{"threshold_breach"})
	if result != nil {
		t.Error("expected nil for non-matching signals")
	}
}

func TestSelectCapsule_PrefersSuccess(t *testing.T) {
	capsules := []Capsule{
		{GeneID: "gene_a", Trigger: []string{"sensor_error"}, Outcome: CapsuleOutcome{Status: "failed"}, Confidence: 0.5},
		{GeneID: "gene_a", Trigger: []string{"sensor_error"}, Outcome: CapsuleOutcome{Status: "success"}, Confidence: 0.5},
	}

	result := SelectCapsule(capsules, "gene_a", []string{"sensor_error"})
	if result == nil {
		t.Fatal("expected a capsule match")
	}
	if result.Outcome.Status != "success" {
		t.Errorf("expected successful capsule to be preferred, got %s", result.Outcome.Status)
	}
}
