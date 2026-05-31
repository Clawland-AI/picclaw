package edge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// GeneStatsProvider is an interface satisfied by the gene.Engine
// to decouple the edge package from pkg/gene.
type GeneStatsProvider interface {
	GetStats() map[string]interface{}
	GetHighConfidenceGenes(minConfidence float64, minVerifiedBy int) []map[string]interface{}
}

// Reporter periodically sends heartbeat and event reports
// to the upstream Fleet Manager (L2 NanoClaw or L3 MoltClaw).
type Reporter struct {
	config       Config
	client       *http.Client
	stopCh       chan struct{}
	stopOnce     sync.Once
	geneProvider GeneStatsProvider
	now          func() time.Time
}

// NewReporter creates a new Edge Reporter.
func NewReporter(cfg Config) *Reporter {
	return &Reporter{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
		stopCh: make(chan struct{}),
		now:    time.Now,
	}
}

// SetHTTPClient overrides the reporter HTTP client. It is primarily useful for
// tests and embedders that need custom transport behavior.
func (r *Reporter) SetHTTPClient(client *http.Client) {
	if client != nil {
		r.client = client
	}
}

// StartHeartbeat begins the periodic heartbeat loop.
func (r *Reporter) StartHeartbeat() {
	interval := r.config.Cloud.HeartbeatInterval
	if interval <= 0 {
		interval = 30 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("[edge] heartbeat reporter started (interval: %s, endpoint: %s)",
		interval, r.config.Cloud.Endpoint)

	// Send initial registration
	if err := r.Register(); err != nil {
		log.Printf("[edge] registration failed: %v", err)
	} else {
		log.Printf("[edge] registered with Fleet Manager as %s", r.config.NodeID)
	}

	for {
		select {
		case <-ticker.C:
			if err := r.SendHeartbeat(); err != nil {
				log.Printf("[edge] heartbeat failed: %v", err)
			}
		case <-r.stopCh:
			log.Println("[edge] heartbeat reporter stopped")
			return
		}
	}
}

// Stop halts the heartbeat reporter.
func (r *Reporter) Stop() {
	r.stopOnce.Do(func() {
		close(r.stopCh)
	})
}

// SetGeneProvider attaches a gene stats provider to include
// gene evolution metrics in heartbeat reports.
func (r *Reporter) SetGeneProvider(gp GeneStatsProvider) {
	r.geneProvider = gp
}

// PublishGene sends a high-confidence gene to the Fleet Manager
// for cross-node sharing via the /fleet/genes/publish endpoint.
func (r *Reporter) PublishGene(geneData map[string]interface{}) error {
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"gene":      geneData,
		"timestamp": r.timestamp(),
	}
	return r.post("/fleet/genes/publish", payload)
}

// ReportEvent sends a one-off event to the Fleet Manager.
func (r *Reporter) ReportEvent(eventType string, data map[string]any) error {
	if strings.TrimSpace(eventType) == "" {
		return fmt.Errorf("event type is required")
	}
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"type":      eventType,
		"data":      data,
		"timestamp": r.timestamp(),
	}
	return r.post("/fleet/events", payload)
}

// Register sends an initial node registration report to the Fleet Manager.
func (r *Reporter) Register() error {
	payload := map[string]any{
		"node_id":      r.config.NodeID,
		"node_name":    r.config.NodeName,
		"type":         "picclaw",
		"capabilities": []string{"sensor", "exec", "cron", "alert"},
		"version":      "0.1.0",
		"timestamp":    r.timestamp(),
	}
	return r.post("/fleet/register", payload)
}

// SendHeartbeat reports that this edge node is online.
func (r *Reporter) SendHeartbeat() error {
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"status":    "online",
		"timestamp": r.timestamp(),
	}

	// Include gene evolution stats if provider is available
	if r.geneProvider != nil {
		payload["gene_stats"] = r.geneProvider.GetStats()
	}

	return r.post("/fleet/heartbeat", payload)
}

// ReportStatus sends a structured status snapshot to the Fleet Manager.
func (r *Reporter) ReportStatus(status string, data map[string]any) error {
	if strings.TrimSpace(status) == "" {
		status = "unknown"
	}
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"node_name": r.config.NodeName,
		"status":    status,
		"data":      data,
		"timestamp": r.timestamp(),
	}
	return r.post("/fleet/status", payload)
}

func (r *Reporter) post(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	url := strings.TrimRight(r.config.Cloud.Endpoint, "/") + path
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if r.config.Cloud.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.config.Cloud.Token)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("post %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		if len(bytes.TrimSpace(detail)) > 0 {
			return fmt.Errorf("post %s: status %d: %s", path, resp.StatusCode, strings.TrimSpace(string(detail)))
		}
		return fmt.Errorf("post %s: status %d", path, resp.StatusCode)
	}
	return nil
}

func (r *Reporter) timestamp() string {
	return r.now().UTC().Format(time.RFC3339)
}
