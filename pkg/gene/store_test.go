package gene

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)
	if err := store.Init(); err != nil {
		t.Fatalf("store.Init() failed: %v", err)
	}
	return store, tmpDir
}

func TestStore_Init_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)
	if err := store.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	genesDir := filepath.Join(tmpDir, "genes")
	if _, err := os.Stat(genesDir); os.IsNotExist(err) {
		t.Error("genes directory was not created")
	}
}

func TestStore_InitEmpty_Counts(t *testing.T) {
	store, _ := setupTestStore(t)
	if store.GeneCount() != 0 {
		t.Errorf("expected 0 genes, got %d", store.GeneCount())
	}
	if store.CapsuleCount() != 0 {
		t.Errorf("expected 0 capsules, got %d", store.CapsuleCount())
	}
}

func TestStore_AddGene_Persists(t *testing.T) {
	store, tmpDir := setupTestStore(t)

	g := Gene{
		ID:           "gene_test",
		Category:     "repair",
		Scenario:     "datacenter",
		SignalsMatch: []string{"sensor_error"},
		Strategy:     []string{"step1"},
		Confidence:   0.8,
	}

	if err := store.AddGene(g); err != nil {
		t.Fatalf("AddGene failed: %v", err)
	}

	if store.GeneCount() != 1 {
		t.Errorf("expected 1 gene, got %d", store.GeneCount())
	}

	// Verify it persisted to disk by creating a new store
	store2 := NewStore(tmpDir)
	if err := store2.Init(); err != nil {
		t.Fatalf("second Init failed: %v", err)
	}
	if store2.GeneCount() != 1 {
		t.Errorf("expected 1 gene after reload, got %d", store2.GeneCount())
	}
}

func TestStore_AddGene_DuplicateIDUpdates(t *testing.T) {
	store, _ := setupTestStore(t)

	g1 := Gene{ID: "gene_dup", Category: "repair", SignalsMatch: []string{"a"}, Strategy: []string{"old"}}
	g2 := Gene{ID: "gene_dup", Category: "optimize", SignalsMatch: []string{"b"}, Strategy: []string{"new"}}

	if err := store.AddGene(g1); err != nil {
		t.Fatalf("AddGene(g1) failed: %v", err)
	}
	if err := store.AddGene(g2); err != nil {
		t.Fatalf("AddGene(g2) failed: %v", err)
	}

	if store.GeneCount() != 1 {
		t.Errorf("expected 1 gene (updated), got %d", store.GeneCount())
	}

	gene := store.GetGene("gene_dup")
	if gene == nil {
		t.Fatal("GetGene returned nil")
	}
	if gene.Category != "optimize" {
		t.Errorf("expected updated category 'optimize', got %s", gene.Category)
	}
}

func TestStore_AddGene_SetsTypeAndVersion(t *testing.T) {
	store, _ := setupTestStore(t)

	g := Gene{ID: "gene_meta", SignalsMatch: []string{"x"}, Strategy: []string{"y"}}
	if err := store.AddGene(g); err != nil {
		t.Fatalf("AddGene failed: %v", err)
	}

	gene := store.GetGene("gene_meta")
	if gene == nil {
		t.Fatal("GetGene returned nil")
	}
	if gene.Type != "Gene" {
		t.Errorf("expected Type 'Gene', got %s", gene.Type)
	}
	if gene.SchemaVersion != "1.0.0" {
		t.Errorf("expected SchemaVersion '1.0.0', got %s", gene.SchemaVersion)
	}
	if gene.AssetID == "" {
		t.Error("expected AssetID to be computed")
	}
}

func TestStore_GetGene_NotFound(t *testing.T) {
	store, _ := setupTestStore(t)
	gene := store.GetGene("nonexistent")
	if gene != nil {
		t.Error("expected nil for non-existent gene")
	}
}

func TestStore_GetGenes_ReturnsCopy(t *testing.T) {
	store, _ := setupTestStore(t)

	g := Gene{ID: "gene_copy", SignalsMatch: []string{"x"}, Strategy: []string{"y"}, Confidence: 0.5}
	if err := store.AddGene(g); err != nil {
		t.Fatalf("AddGene failed: %v", err)
	}

	genes := store.GetGenes()
	genes[0].Confidence = 0.0

	original := store.GetGene("gene_copy")
	if original.Confidence != 0.5 {
		t.Error("GetGenes should return a copy, not a reference")
	}
}

func TestStore_UpdateGene(t *testing.T) {
	store, _ := setupTestStore(t)

	g := Gene{ID: "gene_update", Category: "repair", SignalsMatch: []string{"x"}, Strategy: []string{"y"}, Confidence: 0.5}
	if err := store.AddGene(g); err != nil {
		t.Fatalf("AddGene failed: %v", err)
	}

	g.Confidence = 0.9
	if err := store.UpdateGene(g); err != nil {
		t.Fatalf("UpdateGene failed: %v", err)
	}

	updated := store.GetGene("gene_update")
	if updated.Confidence != 0.9 {
		t.Errorf("expected confidence 0.9, got %f", updated.Confidence)
	}
}

func TestStore_UpdateGene_NotFound(t *testing.T) {
	store, _ := setupTestStore(t)

	g := Gene{ID: "nonexistent", SignalsMatch: []string{"x"}, Strategy: []string{"y"}}
	err := store.UpdateGene(g)
	if err == nil {
		t.Error("expected error when updating non-existent gene")
	}
}

func TestStore_AddCapsule_Persists(t *testing.T) {
	store, tmpDir := setupTestStore(t)

	c := Capsule{
		ID:      "cap_test",
		GeneID:  "gene_a",
		Trigger: []string{"sensor_error"},
		Summary: "test capsule",
		Outcome: CapsuleOutcome{Status: "success"},
	}

	if err := store.AddCapsule(c); err != nil {
		t.Fatalf("AddCapsule failed: %v", err)
	}

	if store.CapsuleCount() != 1 {
		t.Errorf("expected 1 capsule, got %d", store.CapsuleCount())
	}

	// Reload
	store2 := NewStore(tmpDir)
	if err := store2.Init(); err != nil {
		t.Fatalf("reload Init failed: %v", err)
	}
	if store2.CapsuleCount() != 1 {
		t.Errorf("expected 1 capsule after reload, got %d", store2.CapsuleCount())
	}
}

func TestStore_GetCapsulesForGene(t *testing.T) {
	store, _ := setupTestStore(t)

	store.AddCapsule(Capsule{ID: "c1", GeneID: "gene_a", Trigger: []string{"x"}, Summary: "a1", Outcome: CapsuleOutcome{Status: "success"}})
	store.AddCapsule(Capsule{ID: "c2", GeneID: "gene_b", Trigger: []string{"y"}, Summary: "b1", Outcome: CapsuleOutcome{Status: "success"}})
	store.AddCapsule(Capsule{ID: "c3", GeneID: "gene_a", Trigger: []string{"z"}, Summary: "a2", Outcome: CapsuleOutcome{Status: "failed"}})

	capsules := store.GetCapsulesForGene("gene_a")
	if len(capsules) != 2 {
		t.Errorf("expected 2 capsules for gene_a, got %d", len(capsules))
	}
}

func TestStore_SeedGenes_EmptyStore(t *testing.T) {
	store, _ := setupTestStore(t)

	seeds := []Gene{
		{ID: "seed_1", SignalsMatch: []string{"a"}, Strategy: []string{"s1"}},
		{ID: "seed_2", SignalsMatch: []string{"b"}, Strategy: []string{"s2"}},
	}

	if err := store.SeedGenes(seeds); err != nil {
		t.Fatalf("SeedGenes failed: %v", err)
	}

	if store.GeneCount() != 2 {
		t.Errorf("expected 2 genes after seeding, got %d", store.GeneCount())
	}

	gene := store.GetGene("seed_1")
	if gene == nil {
		t.Fatal("seeded gene not found")
	}
	if gene.Type != "Gene" {
		t.Errorf("expected Type 'Gene', got %s", gene.Type)
	}
	if gene.AssetID == "" {
		t.Error("expected AssetID to be computed on seeded gene")
	}
}

func TestStore_SeedGenes_NonEmptyStore_NoOp(t *testing.T) {
	store, _ := setupTestStore(t)

	existing := Gene{ID: "existing", SignalsMatch: []string{"x"}, Strategy: []string{"y"}}
	store.AddGene(existing)

	seeds := []Gene{
		{ID: "seed_1", SignalsMatch: []string{"a"}, Strategy: []string{"s1"}},
	}

	if err := store.SeedGenes(seeds); err != nil {
		t.Fatalf("SeedGenes failed: %v", err)
	}

	if store.GeneCount() != 1 {
		t.Errorf("expected 1 gene (seed should be skipped), got %d", store.GeneCount())
	}
}

func TestStore_AppendEvent(t *testing.T) {
	store, tmpDir := setupTestStore(t)

	evt := EvolutionEvent{
		ID:        "evt_1",
		Intent:    "repair",
		GeneID:    "gene_a",
		CapsuleID: "cap_1",
		Signals:   []string{"sensor_error"},
		Outcome:   CapsuleOutcome{Status: "success"},
		Timestamp: "2026-01-01T00:00:00Z",
	}

	if err := store.AppendEvent(evt); err != nil {
		t.Fatalf("AppendEvent failed: %v", err)
	}

	eventsFile := filepath.Join(tmpDir, "genes", "events.jsonl")
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("failed to read events file: %v", err)
	}
	if len(data) == 0 {
		t.Error("events file should not be empty")
	}
}
