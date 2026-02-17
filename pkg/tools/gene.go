package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sipeed/picoclaw/pkg/gene"
)

// GeneTool allows the agent to explicitly report a learned experience,
// solidifying it as a Gene/Capsule in the evolution system.
type GeneTool struct {
	engine *gene.Engine
}

// NewGeneTool creates a new GeneTool backed by the given Engine.
func NewGeneTool(engine *gene.Engine) *GeneTool {
	return &GeneTool{engine: engine}
}

func (t *GeneTool) Name() string        { return "report_gene" }
func (t *GeneTool) Description() string {
	return "Report a learned experience or monitoring strategy to the Gene Evolution system. " +
		"Use this after successfully handling a novel situation so the experience can be reused. " +
		"The system will create a Gene (reusable strategy) and Capsule (record of this instance)."
}

func (t *GeneTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Brief description of what happened and how it was handled",
			},
			"strategy": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Step-by-step strategy that was used to handle the situation",
			},
			"signals": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Signals or conditions that triggered this situation (e.g. sensor_error, threshold_breach)",
			},
			"outcome": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"success", "failed", "partial"},
				"description": "Result of applying this strategy",
			},
		},
		"required": []string{"summary", "strategy", "signals", "outcome"},
	}
}

func (t *GeneTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	summary, _ := args["summary"].(string)
	outcome, _ := args["outcome"].(string)

	if summary == "" {
		return "", fmt.Errorf("summary is required")
	}
	if outcome == "" {
		outcome = "success"
	}

	// Parse strategy array
	strategy, err := parseStringArray(args, "strategy")
	if err != nil {
		return "", fmt.Errorf("invalid strategy: %w", err)
	}

	// Parse signals array
	signals, err := parseStringArray(args, "signals")
	if err != nil {
		return "", fmt.Errorf("invalid signals: %w", err)
	}

	result, err := t.engine.ReportGene(summary, strategy, signals, outcome)
	if err != nil {
		return "", fmt.Errorf("report gene failed: %w", err)
	}

	// Build response
	response := map[string]interface{}{
		"status":     "solidified",
		"capsule_id": result.CapsuleID,
		"event_id":   result.EventID,
	}

	if result.GeneCreated {
		response["gene_created"] = true
		response["gene_id"] = result.GeneID
		response["message"] = fmt.Sprintf("New gene '%s' created from this experience. It will be available for future similar situations.", result.GeneID)
	} else if result.GeneUpdated {
		response["gene_updated"] = true
		response["gene_id"] = result.GeneID
		response["message"] = fmt.Sprintf("Existing gene '%s' confidence updated based on this outcome.", result.GeneID)
	}

	data, _ := json.MarshalIndent(response, "", "  ")
	return string(data), nil
}

// parseStringArray extracts a []string from args[key].
// Handles both []string and []interface{} (from JSON unmarshaling).
func parseStringArray(args map[string]interface{}, key string) ([]string, error) {
	raw, ok := args[key]
	if !ok {
		return nil, nil
	}

	switch v := raw.(type) {
	case []string:
		return v, nil
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("expected string in array, got %T", item)
			}
			result = append(result, s)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("expected array, got %T", raw)
	}
}
