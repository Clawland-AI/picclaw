package edge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// GeneStatsProvider is satisfied by the gene.Engine.
type GeneStatsProvider interface {
	GetStats() map[string]interface{}
	GetHighConfidenceGenes(minConfidence float64, minVerifiedBy int) []map[string]interface{}
}

// Reporter sends heartbeats, status, and events to the Fleet Manager.
type Reporter struct {
	config Config
	client *http.Client
	stopCh chan struct{}
	closed bool
	mu     sync.Mutex
	wg     sync.WaitGroup
	geneMu sync.RWMutex
	gene   GeneStatsProvider
	now    func() time.Time
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

// Start begins the heartbeat loop with initial registration + heartbeat.
func (r *Reporter) Start() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		hbi := r.config.Cloud.HeartbeatInterval
		if hbi <= 0 {
			hbi = 30 * time.Second
		}
		r.Register()
		r.SendHeartbeat()
		ticker := time.NewTicker(hbi)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.SendHeartbeat()
			case <-r.stopCh:
				return
			}
		}
	}()
}

// StartHeartbeat is a backward-compatible alias for Start().
func (r *Reporter) StartHeartbeat() { r.Start() }

// Stop gracefully shuts down the reporter.
func (r *Reporter) Stop() {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return
	}
	r.closed = true
	r.mu.Unlock()
	close(r.stopCh)
	r.wg.Wait()
}

// SetGeneProvider attaches a gene stats provider for heartbeat enrichment.
func (r *Reporter) SetGeneProvider(gp GeneStatsProvider) {
	r.geneMu.Lock()
	r.gene = gp
	r.geneMu.Unlock()
}

func (r *Reporter) getGeneProvider() GeneStatsProvider {
	r.geneMu.RLock()
	defer r.geneMu.RUnlock()
	return r.gene
}

// Register sends a registration to the Fleet Manager.
func (r *Reporter) Register() error {
	return r.deliver("/fleet/register", map[string]any{
		"node_id":      r.config.NodeID,
		"node_name":    r.config.NodeName,
		"type":         "picclaw",
		"capabilities": []string{"sensor", "exec", "cron", "alert"},
		"version":      "0.2.0",
		"timestamp":    r.now().UTC().Format(time.RFC3339),
	})
}

// SendHeartbeat sends a heartbeat to the Fleet Manager.
func (r *Reporter) SendHeartbeat() error {
	p := map[string]any{
		"node_id":   r.config.NodeID,
		"status":    "online",
		"timestamp": r.now().UTC().Format(time.RFC3339),
	}
	if gp := r.getGeneProvider(); gp != nil {
		p["gene_stats"] = gp.GetStats()
	}
	return r.deliver("/fleet/heartbeat", p)
}

// ReportStatus sends a status snapshot to the Fleet Manager.
func (r *Reporter) ReportStatus() error {
	p := map[string]any{
		"node_id":   r.config.NodeID,
		"node_name": r.config.NodeName,
		"status":    "online",
		"version":   "0.2.0",
		"timestamp": r.now().UTC().Format(time.RFC3339),
	}
	if gp := r.getGeneProvider(); gp != nil {
		p["gene_stats"] = gp.GetStats()
	}
	return r.deliver("/fleet/status", p)
}

// ReportEvent sends an event report to the Fleet Manager.
func (r *Reporter) ReportEvent(eventType string, data map[string]any) error {
	return r.deliver("/fleet/events", map[string]any{
		"node_id":   r.config.NodeID,
		"type":      eventType,
		"data":      data,
		"timestamp": r.now().UTC().Format(time.RFC3339),
	})
}

// deliver sends a JSON POST. Retries once on failure.
func (r *Reporter) deliver(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	err = r.doPost(path, body)
	if err == nil {
		return nil
	}
	time.Sleep(500 * time.Millisecond)
	return r.doPost(path, body)
}

func (r *Reporter) doPost(path string, body []byte) error {
	req, err := http.NewRequest("POST", strings.TrimRight(r.config.Cloud.Endpoint, "/")+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "picclaw-edge-reporter/0.2.0")
	if r.config.NodeID != "" {
		req.Header.Set("X-Clawland-Node-ID", r.config.NodeID)
	}
	if r.config.Cloud.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.config.Cloud.Token)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("post %s: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if s := strings.TrimSpace(string(b)); s != "" {
			return fmt.Errorf("post %s: status %d: %s", path, resp.StatusCode, s)
		}
		return fmt.Errorf("post %s: status %d", path, resp.StatusCode)
	}
	return nil
}
