package gene

import (
	"testing"
)

func setupTestSolidifier(t *testing.T) (*Solidifier, *Store) {
	t.Helper()
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)
	if err := store.Init(); err != nil {
		t.Fatalf("store.Init() failed: %v", err)
	}
	extractor := NewSignalExtractor(tmpDir)
	solidifier := NewSolidifier(store, extractor)
	return solidifier, store
}

func TestSolidify_CreatesNewGeneFromNovelSuccess(t *testing.T) {
	solidifier, store := setupTestSolidifier(t)

	req := SolidifyRequest{
		Signals:  []string{"sensor_error"},
		Strategy: []string{"check sensor", "restart daemon", "verify readings"},
		Summary:  "Fixed sensor read failure by restarting daemon",
		Outcome:  CapsuleOutcome{Status: "success", Details: "sensor recovered"},
	}

	result, err := solidifier.Solidify(req)
	if err != nil {
		t.Fatalf("Solidify failed: %v", err)
	}

	if !result.CapsuleCreated {
		t.Error("expected capsule to be created")
	}
	if result.CapsuleID == "" {
		t.Error("expected non-empty capsule ID")
	}
	if !result.GeneCreated {
		t.Error("expected new gene to be created from novel success")
	}
	if result.GeneID == "" {
		t.Error("expected non-empty gene ID")
	}
	if result.EventID == "" {
		t.Error("expected non-empty event ID")
	}

	if store.GeneCount() != 1 {
		t.Errorf("expected 1 gene, got %d", store.GeneCount())
	}
	if store.CapsuleCount() != 1 {
		t.Errorf("expected 1 capsule, got %d", store.CapsuleCount())
	}
}

func TestSolidify_UpdatesExistingGene(t *testing.T) {
	solidifier, store := setupTestSolidifier(t)

	// Add an existing gene
	store.AddGene(Gene{
		ID:           "gene_existing",
		Category:     "repair",
		SignalsMatch: []string{"sensor_error"},
		Strategy:     []string{"restart"},
		Confidence:   0.7,
		VerifiedBy:   2,
	})

	req := SolidifyRequest{
		Signals:  []string{"sensor_error"},
		Strategy: []string{"restart"},
		Summary:  "Applied existing fix",
		Outcome:  CapsuleOutcome{Status: "success", Details: "fixed again"},
		GeneID:   "gene_existing",
	}

	result, err := solidifier.Solidify(req)
	if err != nil {
		t.Fatalf("Solidify failed: %v", err)
	}

	if !result.GeneUpdated {
		t.Error("expected gene to be updated")
	}
	if result.GeneCreated {
		t.Error("should not create new gene when existing one is specified")
	}

	updated := store.GetGene("gene_existing")
	if updated == nil {
		t.Fatal("gene not found after update")
	}
	if updated.Confidence <= 0.7 {
		t.Errorf("expected confidence increase on success, got %f", updated.Confidence)
	}
	if updated.VerifiedBy <= 2 {
		t.Errorf("expected verification count increase, got %d", updated.VerifiedBy)
	}
}

func TestSolidify_FailedOutcome_DecreasesConfidence(t *testing.T) {
	solidifier, store := setupTestSolidifier(t)

	store.AddGene(Gene{
		ID:           "gene_failing",
		Category:     "repair",
		SignalsMatch: []string{"sensor_error"},
		Strategy:     []string{"restart"},
		Confidence:   0.8,
		VerifiedBy:   3,
	})

	req := SolidifyRequest{
		Signals: []string{"sensor_error"},
		Summary: "Fix failed",
		Outcome: CapsuleOutcome{Status: "failed", Details: "still broken"},
		GeneID:  "gene_failing",
	}

	_, err := solidifier.Solidify(req)
	if err != nil {
		t.Fatalf("Solidify failed: %v", err)
	}

	updated := store.GetGene("gene_failing")
	if updated.Confidence >= 0.8 {
		t.Errorf("expected confidence decrease on failure, got %f", updated.Confidence)
	}
}

func TestSolidify_NoGeneCreatedOnFailure(t *testing.T) {
	solidifier, store := setupTestSolidifier(t)

	req := SolidifyRequest{
		Signals:  []string{"sensor_error"},
		Strategy: []string{"tried something"},
		Summary:  "Novel attempt that failed",
		Outcome:  CapsuleOutcome{Status: "failed", Details: "didn't work"},
	}

	result, err := solidifier.Solidify(req)
	if err != nil {
		t.Fatalf("Solidify failed: %v", err)
	}

	if result.GeneCreated {
		t.Error("should not create gene from failed novel attempt")
	}
	if store.GeneCount() != 0 {
		t.Errorf("expected 0 genes, got %d", store.GeneCount())
	}
	if !result.CapsuleCreated {
		t.Error("capsule should still be created for failed attempt")
	}
}

func TestComputeCapsuleConfidence(t *testing.T) {
	solidifier, _ := setupTestSolidifier(t)

	tests := []struct {
		name     string
		outcome  string
		sensor   map[string]any
		strategy []string
		minConf  float64
		maxConf  float64
	}{
		{"success", "success", nil, nil, 0.7, 0.9},
		{"failed", "failed", nil, nil, 0.1, 0.3},
		{"partial", "partial", nil, nil, 0.4, 0.6},
		{"success+sensor", "success", map[string]any{"temp": 42}, nil, 0.8, 1.0},
		{"success+strategy", "success", nil, []string{"a", "b", "c"}, 0.8, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := SolidifyRequest{
				Outcome:    CapsuleOutcome{Status: tt.outcome},
				SensorData: tt.sensor,
				Strategy:   tt.strategy,
			}
			conf := solidifier.computeCapsuleConfidence(req)
			if conf < tt.minConf || conf > tt.maxConf {
				t.Errorf("expected confidence in [%f, %f], got %f", tt.minConf, tt.maxConf, conf)
			}
		})
	}
}

func TestInferIntent(t *testing.T) {
	solidifier, _ := setupTestSolidifier(t)

	tests := []struct {
		signals  []string
		expected string
	}{
		{[]string{SignalSensorError}, "repair"},
		{[]string{SignalThresholdBreach}, "repair"},
		{[]string{SignalResponseTooSlow}, "optimize"},
		{[]string{SignalStrategyProven}, "optimize"},
		{[]string{SignalTimePattern}, "optimize"},
		{[]string{SignalNewPatternDetected}, "innovate"},
		{[]string{SignalUnknownSituation}, "innovate"},
		{[]string{SignalCrossSensorAnomaly}, "innovate"},
		{[]string{}, "optimize"},
	}

	for _, tt := range tests {
		req := SolidifyRequest{Signals: tt.signals}
		intent := solidifier.inferIntent(req)
		if intent != tt.expected {
			t.Errorf("signals=%v: expected %s, got %s", tt.signals, tt.expected, intent)
		}
	}
}

func TestBuildEnvFingerprint(t *testing.T) {
	solidifier, _ := setupTestSolidifier(t)
	fp := solidifier.buildEnvFingerprint()

	if fp.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
	if fp.Platform == "" {
		t.Error("Platform should not be empty")
	}
	if fp.Arch == "" {
		t.Error("Arch should not be empty")
	}
}
