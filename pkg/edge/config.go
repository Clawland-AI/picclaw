// Package edge provides the Edge API Server and Reporter for PicoClaw.
// This enables PicoClaw to act as an L1 edge node in the Clawland
// cloud-edge architecture, communicating with L2 (NanoClaw) and L3 (MoltClaw).
package edge

import "time"

// Config holds the Edge API Server configuration.
type Config struct {
	// Enabled activates the Edge API Server.
	Enabled bool `json:"enabled"`

	// Port is the HTTP port for the Edge API Server.
	Port int `json:"port"`

	// NodeID is the unique identifier for this edge node.
	NodeID string `json:"node_id"`

	// NodeName is the human-readable name for this edge node.
	NodeName string `json:"node_name"`

	// Cloud holds the upstream (L2/L3) connection settings.
	Cloud CloudConfig `json:"cloud"`

	// StatusInterval is how often status snapshots are sent to the Fleet Manager.
	// Defaults to 5 minutes if zero.
	StatusInterval time.Duration `json:"status_interval"`

	// MaxRetries is the maximum number of delivery retries for failed HTTP calls.
	// Defaults to 3 if zero.
	MaxRetries int `json:"max_retries"`

	// RetryBaseDelay is the initial backoff duration for retries.
	// Each retry doubles this with jitter. Defaults to 1s if zero.
	RetryBaseDelay time.Duration `json:"retry_base_delay"`

	// EventQueueSize is the capacity of the buffered event queue for async reporting.
	// Defaults to 256 if zero.
	EventQueueSize int `json:"event_queue_size"`

	// EventWorkers is the number of goroutines draining the async event queue.
	// Defaults to 1 if zero.
	EventWorkers int `json:"event_workers"`
}

// CloudConfig holds the upstream cloud/gateway connection settings.
type CloudConfig struct {
	// Endpoint is the URL of the Fleet Manager (L2 NanoClaw or L3 MoltClaw).
	Endpoint string `json:"endpoint"`

	// HeartbeatInterval is the interval between heartbeat reports.
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`

	// Token is the authentication token for the Fleet Manager.
	Token string `json:"token,omitempty"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:  false,
		Port:     9090,
		NodeID:   "",
		NodeName: "picclaw-edge",
		Cloud: CloudConfig{
			Endpoint:          "http://localhost:8080",
			HeartbeatInterval: 30 * time.Second,
		},
		StatusInterval:  5 * time.Minute,
		MaxRetries:      3,
		RetryBaseDelay:  1 * time.Second,
		EventQueueSize:  256,
		EventWorkers:    1,
	}
}
