// Clawland Gene Evolution Protocol (CGEP)
// Adapted from autogame-17/evolver GEP protocol for physical-world edge AI.
// Genes encode monitoring strategies, not code patches.

package gene

// Gene represents a reusable monitoring strategy pattern.
// Adapted from Evolver's genes.json schema for physical-world scenarios.
type Gene struct {
	Type           string          `json:"type"`            // always "Gene"
	ID             string          `json:"id"`              // e.g. "gene_dc_temp_spike_crosscheck"
	SchemaVersion  string          `json:"schema_version"`  // "1.0.0"
	AssetID        string          `json:"asset_id"`        // SHA-256 content hash
	Category       string          `json:"category"`        // "repair" | "optimize" | "innovate"
	Scenario       string          `json:"scenario"`        // "datacenter" | "aquaculture" | "greenhouse" | "generic"
	SignalsMatch   []string        `json:"signals_match"`   // trigger conditions
	Strategy       []string        `json:"strategy"`        // step-by-step handling instructions
	Constraints    GeneConstraints `json:"constraints"`     // safety constraints
	Validation     []ValidationStep `json:"validation"`     // how to verify this gene works
	Confidence     float64         `json:"confidence"`      // 0.0 - 1.0
	VerifiedBy     int             `json:"verified_by"`     // nodes that confirmed effectiveness
	EnvFingerprint EnvFingerprint  `json:"env_fingerprint"` // hardware/software context
	ParentGene     string          `json:"parent_gene"`     // mutated from which gene (empty if original)
	CreatedAt      string          `json:"created_at"`
}

// GeneConstraints defines safety boundaries for a Gene.
type GeneConstraints struct {
	MaxRetries     int      `json:"max_retries"`     // max retry attempts
	TimeoutSeconds int      `json:"timeout_seconds"` // max time for strategy execution
	Requires       []string `json:"requires"`        // required sensors/capabilities
	Forbids        []string `json:"forbids"`         // forbidden actions
}

// ValidationStep defines how to verify a Gene's effectiveness.
type ValidationStep struct {
	Command     string `json:"command"`     // validation command or check
	Description string `json:"description"` // what this step validates
	Expected    string `json:"expected"`    // expected outcome
}

// Capsule is a solidified instance of a Gene application.
// Records what happened when a Gene was applied in a specific context.
type Capsule struct {
	Type           string         `json:"type"`            // always "Capsule"
	ID             string         `json:"id"`
	SchemaVersion  string         `json:"schema_version"`
	AssetID        string         `json:"asset_id"`
	GeneID         string         `json:"gene_id"`         // which Gene was used
	Trigger        []string       `json:"trigger"`         // actual signals that triggered this
	Summary        string         `json:"summary"`         // what happened
	Confidence     float64        `json:"confidence"`
	Outcome        CapsuleOutcome `json:"outcome"`         // success/failed + details
	EnvFingerprint EnvFingerprint `json:"env_fingerprint"`
	SensorData     map[string]any `json:"sensor_data"`     // actual readings at time of event
	CreatedAt      string         `json:"created_at"`
}

// CapsuleOutcome records the result of applying a Gene.
type CapsuleOutcome struct {
	Status  string         `json:"status"`  // "success" | "failed" | "partial"
	Details string         `json:"details"` // human-readable description
	Metrics map[string]any `json:"metrics"` // measurable results (response_time_ms, etc.)
}

// EnvFingerprint captures the hardware/software context.
// Critical for edge devices where behavior varies by hardware.
type EnvFingerprint struct {
	Hardware    string   `json:"hardware"`     // "LicheeRV-Nano" | "MaixCAM" | ...
	GoVersion   string   `json:"go_version"`
	Platform    string   `json:"platform"`     // "linux" | "windows"
	Arch        string   `json:"arch"`         // "riscv64" | "arm64" | "amd64"
	Sensors     []string `json:"sensors"`      // ["DHT22", "DO-Sensor"]
	LLMProvider string   `json:"llm_provider"` // "zhipu" | "groq" | "openrouter"
}

// EvolutionEvent is the audit trail entry for Gene evolution.
type EvolutionEvent struct {
	Type      string         `json:"type"`       // always "EvolutionEvent"
	ID        string         `json:"id"`         // "evt_1234567890"
	Intent    string         `json:"intent"`     // "repair" | "optimize" | "innovate"
	GeneID    string         `json:"gene_id"`
	CapsuleID string         `json:"capsule_id"`
	Signals   []string       `json:"signals"`
	Outcome   CapsuleOutcome `json:"outcome"`
	Timestamp string         `json:"timestamp"`
}

// StrategyPreset defines Gene application thresholds.
type StrategyPreset struct {
	Name               string  `json:"name"`                 // "conservative" | "balanced" | "exploratory" | "repair-only"
	MinConfidence      float64 `json:"min_confidence"`       // minimum confidence to apply
	MinVerifiedBy      int     `json:"min_verified_by"`      // minimum verifications to apply
	AllowUnverified    bool    `json:"allow_unverified"`     // allow genes with 0 verifications
	CategoryFilter     string  `json:"category_filter"`      // empty = all, or "repair" only
}

// GeneConfig holds runtime configuration for the Gene Engine.
type GeneConfig struct {
	Strategy       string `json:"strategy"`         // preset name
	AutoPublish    bool   `json:"auto_publish"`     // publish verified genes to fleet
	MinConfidence  float64 `json:"min_confidence"`  // override preset threshold
	MinVerifiedBy  int    `json:"min_verified_by"`  // override preset threshold
}

// DefaultGeneConfig returns the default balanced configuration.
func DefaultGeneConfig() GeneConfig {
	return GeneConfig{
		Strategy:      "balanced",
		AutoPublish:   true,
		MinConfidence: 0.7,
		MinVerifiedBy: 2,
	}
}

// Predefined strategy presets adapted from Evolver's EVOLVE_STRATEGY.
var StrategyPresets = map[string]StrategyPreset{
	"conservative": {
		Name:            "conservative",
		MinConfidence:   0.9,
		MinVerifiedBy:   5,
		AllowUnverified: false,
		CategoryFilter:  "",
	},
	"balanced": {
		Name:            "balanced",
		MinConfidence:   0.7,
		MinVerifiedBy:   2,
		AllowUnverified: false,
		CategoryFilter:  "",
	},
	"exploratory": {
		Name:            "exploratory",
		MinConfidence:   0.3,
		MinVerifiedBy:   0,
		AllowUnverified: true,
		CategoryFilter:  "",
	},
	"repair-only": {
		Name:            "repair-only",
		MinConfidence:   0.5,
		MinVerifiedBy:   1,
		AllowUnverified: false,
		CategoryFilter:  "repair",
	},
}
