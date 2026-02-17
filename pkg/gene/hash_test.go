package gene

import (
	"testing"
)

func TestComputeGeneAssetID_Deterministic(t *testing.T) {
	g := Gene{
		ID:           "gene_test_1",
		Category:     "repair",
		Scenario:     "datacenter",
		SignalsMatch: []string{"sensor_error", "threshold_breach"},
		Strategy:     []string{"step1", "step2"},
		Constraints:  GeneConstraints{MaxRetries: 3, TimeoutSeconds: 300},
	}

	id1 := ComputeGeneAssetID(&g)
	id2 := ComputeGeneAssetID(&g)

	if id1 != id2 {
		t.Errorf("expected deterministic hash, got %s and %s", id1, id2)
	}
	if id1 == "" {
		t.Error("asset ID should not be empty")
	}
	if len(id1) < 10 {
		t.Errorf("asset ID suspiciously short: %s", id1)
	}
}

func TestComputeGeneAssetID_DifferentContent(t *testing.T) {
	g1 := Gene{
		ID:           "gene_a",
		Category:     "repair",
		Scenario:     "datacenter",
		SignalsMatch: []string{"sensor_error"},
		Strategy:     []string{"fix it"},
	}
	g2 := Gene{
		ID:           "gene_b",
		Category:     "optimize",
		Scenario:     "aquaculture",
		SignalsMatch: []string{"threshold_breach"},
		Strategy:     []string{"tune it"},
	}

	id1 := ComputeGeneAssetID(&g1)
	id2 := ComputeGeneAssetID(&g2)

	if id1 == id2 {
		t.Errorf("different genes should have different asset IDs")
	}
}

func TestComputeGeneAssetID_SignalOrderIndependent(t *testing.T) {
	g1 := Gene{
		ID:           "gene_order",
		Category:     "repair",
		Scenario:     "datacenter",
		SignalsMatch: []string{"alpha", "beta", "gamma"},
		Strategy:     []string{"step"},
	}
	g2 := Gene{
		ID:           "gene_order",
		Category:     "repair",
		Scenario:     "datacenter",
		SignalsMatch: []string{"gamma", "alpha", "beta"},
		Strategy:     []string{"step"},
	}

	id1 := ComputeGeneAssetID(&g1)
	id2 := ComputeGeneAssetID(&g2)

	if id1 != id2 {
		t.Errorf("signal order should not affect asset ID; got %s and %s", id1, id2)
	}
}

func TestComputeGeneAssetID_HasPrefix(t *testing.T) {
	g := Gene{
		ID:           "gene_prefix",
		Category:     "repair",
		SignalsMatch: []string{"x"},
		Strategy:     []string{"y"},
	}
	id := ComputeGeneAssetID(&g)
	if len(id) < 7 || id[:7] != "sha256-" {
		t.Errorf("asset ID should start with 'sha256-', got %s", id)
	}
}

func TestComputeCapsuleAssetID_Deterministic(t *testing.T) {
	c := Capsule{
		GeneID:  "gene_test",
		Trigger: []string{"sensor_error"},
		Summary: "test capsule",
		Outcome: CapsuleOutcome{Status: "success", Details: "all good"},
	}

	id1 := ComputeCapsuleAssetID(&c)
	id2 := ComputeCapsuleAssetID(&c)

	if id1 != id2 {
		t.Errorf("expected deterministic hash, got %s and %s", id1, id2)
	}
}

func TestComputeCapsuleAssetID_DifferentContent(t *testing.T) {
	c1 := Capsule{
		GeneID:  "gene_a",
		Trigger: []string{"sensor_error"},
		Summary: "first capsule",
		Outcome: CapsuleOutcome{Status: "success"},
	}
	c2 := Capsule{
		GeneID:  "gene_b",
		Trigger: []string{"threshold_breach"},
		Summary: "second capsule",
		Outcome: CapsuleOutcome{Status: "failed"},
	}

	id1 := ComputeCapsuleAssetID(&c1)
	id2 := ComputeCapsuleAssetID(&c2)

	if id1 == id2 {
		t.Errorf("different capsules should have different asset IDs")
	}
}

func TestSortedCopy(t *testing.T) {
	input := []string{"charlie", "alpha", "bravo"}
	sorted := sortedCopy(input)

	if sorted[0] != "alpha" || sorted[1] != "bravo" || sorted[2] != "charlie" {
		t.Errorf("expected sorted slice, got %v", sorted)
	}

	// Original should not be modified
	if input[0] != "charlie" {
		t.Error("sortedCopy should not modify original slice")
	}
}
