package edge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	config        Config
	client        *http.Client
	stopCh        chan struct{}
	geneProvider  GeneStatsProvider
}

// NewReporter creates a new Edge Reporter.
func NewReporter(cfg Config) *Reporter {
	return &Reporter{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
		stopCh: make(chan struct{}),
	}
}

// StartHeartbeat begins the periodic heartbeat loop.
func (r *Reporter) StartHeartbeat() {
	ticker := time.NewTicker(r.config.Cloud.HeartbeatInterval)
	defer ticker.Stop()

	log.Printf("[edge] heartbeat reporter started (interval: %s, endpoint: %s)",
		r.config.Cloud.HeartbeatInterval, r.config.Cloud.Endpoint)

	// Send initial registration
	r.sendRegistration()

	for {
		select {
		case <-ticker.C:
			r.sendHeartbeat()
		case <-r.stopCh:
			log.Println("[edge] heartbeat reporter stopped")
			return
		}
	}
}

// Stop halts the heartbeat reporter.
func (r *Reporter) Stop() {
	close(r.stopCh)
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
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	return r.post("/fleet/genes/publish", payload)
}

// ReportEvent sends a one-off event to the Fleet Manager.
func (r *Reporter) ReportEvent(eventType string, data map[string]any) error {
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"type":      eventType,
		"data":      data,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	return r.post("/fleet/events", payload)
}

func (r *Reporter) sendRegistration() {
	payload := map[string]any{
		"node_id":      r.config.NodeID,
		"node_name":    r.config.NodeName,
		"type":         "picclaw",
		"capabilities": []string{"sensor", "exec", "cron", "alert"},
		"version":      "0.1.0",
	}
	if err := r.post("/fleet/register", payload); err != nil {
		log.Printf("[edge] registration failed: %v", err)
	} else {
		log.Printf("[edge] registered with Fleet Manager as %s", r.config.NodeID)
	}
}

func (r *Reporter) sendHeartbeat() {
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"status":    "online",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	// Include gene evolution stats if provider is available
	if r.geneProvider != nil {
		payload["gene_stats"] = r.geneProvider.GetStats()
	}

	if err := r.post("/fleet/heartbeat", payload); err != nil {
		log.Printf("[edge] heartbeat failed: %v", err)
	}
}

func (r *Reporter) post(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	url := r.config.Cloud.Endpoint + path
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
		return fmt.Errorf("post %s: status %d", path, resp.StatusCode)
	}
	return nil
}