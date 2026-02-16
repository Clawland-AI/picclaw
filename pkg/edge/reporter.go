package edge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Reporter periodically sends heartbeat and event reports
// to the upstream Fleet Manager (L2 NanoClaw or L3 MoltClaw).
type Reporter struct {
	config   Config
	client   *http.Client
	stopCh   chan struct{}
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