package edge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultHeartbeatInterval = 30 * time.Second
	defaultStatusInterval    = 60 * time.Second
	defaultEventBufferSize   = 64
)

// GeneStatsProvider is an interface satisfied by the gene.Engine
// to decouple the edge package from pkg/gene.
type GeneStatsProvider interface {
	GetStats() map[string]interface{}
	GetHighConfidenceGenes(minConfidence float64, minVerifiedBy int) []map[string]interface{}
}

// ReporterEvent is queued when edge subsystems need to publish an
// event-driven report without blocking their own hot path.
type ReporterEvent struct {
	Type      string         `json:"type"`
	Data      map[string]any `json:"data,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// Reporter periodically sends heartbeat and event reports
// to the upstream Fleet Manager (L2 NanoClaw or L3 MoltClaw).
type Reporter struct {
	config       Config
	client       *http.Client
	stopCh       chan struct{}
	stopOnce     sync.Once
	eventCh      chan ReporterEvent
	started      time.Time
	geneProvider GeneStatsProvider
}

// NewReporter creates a new Edge Reporter.
func NewReporter(cfg Config) *Reporter {
	normalizeReporterConfig(&cfg)
	return &Reporter{
		config:  cfg,
		client:  &http.Client{Timeout: 10 * time.Second},
		stopCh:  make(chan struct{}),
		eventCh: make(chan ReporterEvent, defaultEventBufferSize),
		started: time.Now(),
	}
}

// StartHeartbeat begins the periodic heartbeat loop.
func (r *Reporter) StartHeartbeat() {
	heartbeatTicker := time.NewTicker(r.config.Cloud.HeartbeatInterval)
	statusTicker := time.NewTicker(r.config.Cloud.StatusInterval)
	defer heartbeatTicker.Stop()
	defer statusTicker.Stop()

	log.Printf("[edge] reporter started (heartbeat: %s, status: %s, endpoint: %s)",
		r.config.Cloud.HeartbeatInterval, r.config.Cloud.StatusInterval, r.config.Cloud.Endpoint)

	r.sendRegistration()
	r.sendHeartbeat()
	r.sendStatus()

	for {
		select {
		case <-heartbeatTicker.C:
			r.sendHeartbeat()
		case <-statusTicker.C:
			r.sendStatus()
		case event := <-r.eventCh:
			if err := r.postReporterEvent(event); err != nil {
				log.Printf("[edge] event report failed: %v", err)
			}
		case <-r.stopCh:
			log.Println("[edge] reporter stopped")
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

// SetHTTPClient overrides the reporter HTTP client. It is useful for tests
// and embedded deployments that need custom transports or timeouts.
func (r *Reporter) SetHTTPClient(client *http.Client) {
	if client != nil {
		r.client = client
	}
}

// SetGeneProvider attaches a gene stats provider to include
// gene evolution metrics in heartbeat reports.
func (r *Reporter) SetGeneProvider(gp GeneStatsProvider) {
	r.geneProvider = gp
}

// QueueEvent schedules an event for asynchronous delivery by the reporter loop.
// It returns false when the bounded queue is full so callers can decide whether
// to retry or drop low-priority telemetry.
func (r *Reporter) QueueEvent(eventType string, data map[string]any) bool {
	event := ReporterEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	select {
	case r.eventCh <- event:
		return true
	default:
		log.Printf("[edge] event queue full; dropping event type=%s", eventType)
		return false
	}
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
	return r.postReporterEvent(ReporterEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// ReportStatus sends a full status report immediately.
func (r *Reporter) ReportStatus(status string, data map[string]any) error {
	payload := r.statusPayload(status)
	if data != nil {
		payload["data"] = data
	}
	return r.post("/fleet/status", payload)
}

func (r *Reporter) postReporterEvent(event ReporterEvent) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"type":      event.Type,
		"data":      event.Data,
		"timestamp": event.Timestamp.Format(time.RFC3339),
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
		"node_id":        r.config.NodeID,
		"node_name":      r.config.NodeName,
		"status":         "online",
		"uptime_seconds": int(time.Since(r.started).Seconds()),
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	}

	if r.geneProvider != nil {
		payload["gene_stats"] = r.geneProvider.GetStats()
	}

	if err := r.post("/fleet/heartbeat", payload); err != nil {
		log.Printf("[edge] heartbeat failed: %v", err)
	}
}

func (r *Reporter) sendStatus() {
	if err := r.ReportStatus("online", nil); err != nil {
		log.Printf("[edge] status report failed: %v", err)
	}
}

func (r *Reporter) statusPayload(status string) map[string]any {
	payload := map[string]any{
		"node_id":        r.config.NodeID,
		"node_name":      r.config.NodeName,
		"type":           "picclaw",
		"status":         status,
		"cloud_endpoint": r.config.Cloud.Endpoint,
		"uptime_seconds": int(time.Since(r.started).Seconds()),
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	}
	if r.geneProvider != nil {
		payload["gene_stats"] = r.geneProvider.GetStats()
		payload["high_confidence_genes"] = r.geneProvider.GetHighConfidenceGenes(0.7, 3)
	}
	return payload
}

func (r *Reporter) post(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	endpoint := strings.TrimRight(r.config.Cloud.Endpoint, "/")
	if endpoint == "" {
		return fmt.Errorf("cloud endpoint is empty")
	}

	url := endpoint + path
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Clawland-Node-ID", r.config.NodeID)
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

func normalizeReporterConfig(cfg *Config) {
	if cfg.Cloud.HeartbeatInterval <= 0 {
		cfg.Cloud.HeartbeatInterval = defaultHeartbeatInterval
	}
	if cfg.Cloud.StatusInterval <= 0 {
		cfg.Cloud.StatusInterval = defaultStatusInterval
	}
}
