package gene

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store manages local Gene/Capsule storage as JSON files.
// Storage layout:
//
//	workspace/genes/genes.json      - Gene definitions
//	workspace/genes/capsules.json   - Capsule instances
//	workspace/genes/events.jsonl    - Append-only evolution audit trail
type Store struct {
	dir          string
	genesFile    string
	capsulesFile string
	eventsFile   string

	genes    []Gene
	capsules []Capsule
	mu       sync.RWMutex
}

// NewStore creates a new Store that persists to workspace/genes/.
func NewStore(workspace string) *Store {
	dir := filepath.Join(workspace, "genes")
	return &Store{
		dir:          dir,
		genesFile:    filepath.Join(dir, "genes.json"),
		capsulesFile: filepath.Join(dir, "capsules.json"),
		eventsFile:   filepath.Join(dir, "events.jsonl"),
	}
}

// Init ensures the genes directory exists and loads existing data.
// If no genes.json exists, it will be created with seed genes on first save.
func (s *Store) Init() error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("create genes dir: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Load genes
	if data, err := os.ReadFile(s.genesFile); err == nil {
		if err := json.Unmarshal(data, &s.genes); err != nil {
			return fmt.Errorf("parse genes.json: %w", err)
		}
	}

	// Load capsules
	if data, err := os.ReadFile(s.capsulesFile); err == nil {
		if err := json.Unmarshal(data, &s.capsules); err != nil {
			return fmt.Errorf("parse capsules.json: %w", err)
		}
	}

	return nil
}

// GetGenes returns a copy of all genes.
func (s *Store) GetGenes() []Gene {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Gene, len(s.genes))
	copy(result, s.genes)
	return result
}

// GetGene returns a gene by ID, or nil if not found.
func (s *Store) GetGene(id string) *Gene {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.genes {
		if s.genes[i].ID == id {
			g := s.genes[i]
			return &g
		}
	}
	return nil
}

// GetCapsules returns a copy of all capsules.
func (s *Store) GetCapsules() []Capsule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Capsule, len(s.capsules))
	copy(result, s.capsules)
	return result
}

// GetCapsulesForGene returns capsules that reference a specific gene.
func (s *Store) GetCapsulesForGene(geneID string) []Capsule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Capsule
	for _, c := range s.capsules {
		if c.GeneID == geneID {
			result = append(result, c)
		}
	}
	return result
}

// AddGene adds a new gene (computes asset ID) and persists to disk.
func (s *Store) AddGene(g Gene) error {
	g.Type = "Gene"
	g.SchemaVersion = "1.0.0"
	g.AssetID = ComputeGeneAssetID(&g)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate ID
	for i, existing := range s.genes {
		if existing.ID == g.ID {
			s.genes[i] = g
			return s.saveGenesLocked()
		}
	}

	s.genes = append(s.genes, g)
	return s.saveGenesLocked()
}

// UpdateGene updates an existing gene and persists.
func (s *Store) UpdateGene(g Gene) error {
	g.AssetID = ComputeGeneAssetID(&g)

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.genes {
		if existing.ID == g.ID {
			s.genes[i] = g
			return s.saveGenesLocked()
		}
	}
	return fmt.Errorf("gene %s not found", g.ID)
}

// AddCapsule adds a new capsule (computes asset ID) and persists.
func (s *Store) AddCapsule(c Capsule) error {
	c.Type = "Capsule"
	c.SchemaVersion = "1.0.0"
	c.AssetID = ComputeCapsuleAssetID(&c)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.capsules = append(s.capsules, c)
	return s.saveCapsulesLocked()
}

// AppendEvent appends an evolution event to the JSONL audit trail.
func (s *Store) AppendEvent(evt EvolutionEvent) error {
	evt.Type = "EvolutionEvent"

	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(s.eventsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open events file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write event: %w", err)
	}
	return nil
}

// GeneCount returns the number of local genes.
func (s *Store) GeneCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.genes)
}

// CapsuleCount returns the number of local capsules.
func (s *Store) CapsuleCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.capsules)
}

// SeedGenes initializes the store with seed genes if it's currently empty.
func (s *Store) SeedGenes(seeds []Gene) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.genes) > 0 {
		return nil // already has genes, don't overwrite
	}

	for i := range seeds {
		seeds[i].Type = "Gene"
		seeds[i].SchemaVersion = "1.0.0"
		seeds[i].AssetID = ComputeGeneAssetID(&seeds[i])
	}

	s.genes = seeds
	return s.saveGenesLocked()
}

// saveGenesLocked persists genes to disk. Caller must hold s.mu.
func (s *Store) saveGenesLocked() error {
	data, err := json.MarshalIndent(s.genes, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal genes: %w", err)
	}
	return os.WriteFile(s.genesFile, data, 0644)
}

// saveCapsulesLocked persists capsules to disk. Caller must hold s.mu.
func (s *Store) saveCapsulesLocked() error {
	data, err := json.MarshalIndent(s.capsules, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal capsules: %w", err)
	}
	return os.WriteFile(s.capsulesFile, data, 0644)
}
