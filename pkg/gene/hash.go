package gene

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ComputeAssetID computes a SHA-256 content-addressable asset ID for a Gene.
// Deterministic: same content always produces the same hash.
// Ported from Evolver's contentHash.js concept.
func ComputeGeneAssetID(g *Gene) string {
	// Hash the core content fields (not metadata like created_at or asset_id)
	parts := []string{
		g.ID,
		g.Category,
		g.Scenario,
		strings.Join(sortedCopy(g.SignalsMatch), "|"),
		strings.Join(g.Strategy, "|"),
	}

	constraintsJSON, _ := json.Marshal(g.Constraints)
	parts = append(parts, string(constraintsJSON))

	content := strings.Join(parts, "\n")
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("sha256-%x", hash)
}

// ComputeCapsuleAssetID computes a SHA-256 content-addressable asset ID for a Capsule.
func ComputeCapsuleAssetID(c *Capsule) string {
	parts := []string{
		c.GeneID,
		strings.Join(sortedCopy(c.Trigger), "|"),
		c.Summary,
		c.Outcome.Status,
		c.Outcome.Details,
	}

	sensorJSON, _ := json.Marshal(c.SensorData)
	parts = append(parts, string(sensorJSON))

	content := strings.Join(parts, "\n")
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("sha256-%x", hash)
}

// sortedCopy returns a sorted copy of a string slice (for deterministic hashing).
func sortedCopy(s []string) []string {
	c := make([]string, len(s))
	copy(c, s)
	sort.Strings(c)
	return c
}
