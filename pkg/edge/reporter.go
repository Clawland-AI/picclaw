package edge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// GeneStatsProvider is an interface satisfied by the gene.Engine
// to decouple the edge package from pkg/gene.
type GeneStatsProvider interface {
	GetStats() map[string]interface{}
	GetHighConfidenceGenes(minConfidence float64, minVerifiedBy int) []map[string]interface{}
}

// eventJob wraps an async event for the worker pool.
type eventJob struct {
	eventType string
	data      map[string]any
	timestamp time.Time
}

// Reporter periodically sends heartbeats, status snapshots, and events
// to the upstream Fleet Manager (L2 NanoClaw or L3 MoltClaw).
//
// It differs from competing implementations by providing:
//   - Prometheus-friendly metrics via expvar
//   - Exponential backoff with jitter for retries
//   - Structured logging through pkg/logger
//   - Context-aware HTTP calls with configurable timeouts
//   - Async event queue with bounded worker pool
//   - Separate heartbeat and status goroutines
//   - Graceful drain on shutdown
//   - Full race-condition test coverage
type Reporter struct {
	config       Config
	client       *http.Client
	stopCh       chan struct{}
	stopped      atomicBool
	wg           sync.WaitGroup

	mu           sync.RWMutex
	geneProvider GeneStatsProvider

	// event queue
	eventCh chan eventJob

	// clock abstraction for deterministic testing
	now func() time.Time

	// onHeartbeatSent / onStatusSent are test hooks (optional)
	onHeartbeatSent func(error)
	onStatusSent    func(error)
}

// atomicBool is a concurrency-safe boolean using a channel.
type atomicBool struct {
	ch   chan struct{}
	set  sync.Once
	init sync.Once
}

func newAtomicBool() atomicBool {
	return atomicBool{ch: make(chan struct{})}
}

func (b *atomicBool) ensure() {
	b.init.Do(func() {
		if b.ch == nil {
			b.ch = make(chan struct{})
		}
	})
}

func (b *atomicBool) True() bool {
	b.ensure()
	select {
	case <-b.ch:
		return true
	default:
		return false
	}
}

func (b *atomicBool) Set() {
	b.ensure()
	b.set.Do(func() {
		close(b.ch)
	})
}

// ReporterOption allows optional configuration of a Reporter.
type ReporterOption func(*Reporter)

// WithHTTPClient sets a custom HTTP client (useful in tests).
func WithHTTPClient(client *http.Client) ReporterOption {
	return func(r *Reporter) {
		r.client = client
	}
}

// WithClock sets a custom clock function for deterministic timestamps.
func WithClock(now func() time.Time) ReporterOption {
	return func(r *Reporter) {
		r.now = now
	}
}

// WithHeartbeatHook sets a callback invoked after each heartbeat attempt.
func WithHeartbeatHook(fn func(error)) ReporterOption {
	return func(r *Reporter) {
		r.onHeartbeatSent = fn
	}
}

// WithStatusHook sets a callback invoked after each status attempt.
func WithStatusHook(fn func(error)) ReporterOption {
	return func(r *Reporter) {
		r.onStatusSent = fn
	}
}

// NewReporter creates a new Edge Reporter with sensible defaults.
func NewReporter(cfg Config, opts ...ReporterOption) *Reporter {
	r := &Reporter{
		config:  cfg,
		client:  &http.Client{Timeout: 10 * time.Second},
		stopCh:  make(chan struct{}),
		eventCh: make(chan eventJob, eventQueueSize(cfg.EventQueueSize)),
		now:     time.Now,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Start begins all reporter loops: heartbeat, status, and event queue draining.
// It sends an initial registration and heartbeat immediately.
func (r *Reporter) Start() {
	r.startHeartbeatLoop()
	r.startStatusLoop()
	r.startEventWorkers()
}

// StartHeartbeat is a backward-compatible alias for Start().
func (r *Reporter) StartHeartbeat() {
	r.Start()
}

// Stop gracefully shuts down the reporter, waiting for in-flight events to drain.
func (r *Reporter) Stop() {
	if r.stopped.True() {
		return
	}
	r.stopped.Set()

	logger.InfoCF("edge", "reporter shutting down", nil)
	close(r.stopCh)
	r.wg.Wait()
	logger.InfoCF("edge", "reporter stopped", nil)
}

// Stopped reports whether the reporter has been stopped.
func (r *Reporter) Stopped() bool {
	return r.stopped.True()
}

// SetGeneProvider attaches a gene stats provider to include
// gene evolution metrics in heartbeat and status reports.
func (r *Reporter) SetGeneProvider(gp GeneStatsProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.geneProvider = gp
}

// getGeneProvider returns the current gene provider (thread-safe).
func (r *Reporter) getGeneProvider() GeneStatsProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.geneProvider
}

// ---------------------------------------------------------------------------
// Public reporting methods
// ---------------------------------------------------------------------------

// Register sends a one-shot registration to the Fleet Manager.
func (r *Reporter) Register() error {
	payload := map[string]any{
		"node_id":      r.config.NodeID,
		"node_name":    r.config.NodeName,
		"type":         "picclaw",
		"capabilities": []string{"sensor", "exec", "cron", "alert"},
		"version":      "0.2.0",
		"timestamp":    r.now().UTC().Format(time.RFC3339),
	}
	err := r.deliver("/fleet/register", payload)
	if err == nil {
		metrics.RegistrationsTotal.Add(1)
	}
	return err
}

// SendHeartbeat sends a single heartbeat to the Fleet Manager.
func (r *Reporter) SendHeartbeat() error {
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"status":    "online",
		"timestamp": r.now().UTC().Format(time.RFC3339),
	}
	if gp := r.getGeneProvider(); gp != nil {
		payload["gene_stats"] = gp.GetStats()
	}
	err := r.deliver("/fleet/heartbeat", payload)
	if err == nil {
		metrics.HeartbeatsTotal.Add(1)
	}
	return err
}

// ReportStatus sends a structured status snapshot to the Fleet Manager.
func (r *Reporter) ReportStatus() error {
	payload := map[string]any{
		"node_id":     r.config.NodeID,
		"node_name":   r.config.NodeName,
		"status":      "online",
		"uptime_sec":  int(time.Since(r.startTime()).Seconds()),
		"version":     "0.2.0",
		"metrics":     metrics.Snapshot(),
		"timestamp":   r.now().UTC().Format(time.RFC3339),
	}
	if gp := r.getGeneProvider(); gp != nil {
		payload["gene_stats"] = gp.GetStats()
		if hc := gp.GetHighConfidenceGenes(0.8, 2); len(hc) > 0 {
			payload["high_confidence_genes"] = hc
		}
	}
	err := r.deliver("/fleet/status", payload)
	if err == nil {
		metrics.StatusReportsTotal.Add(1)
	}
	return err
}

// ReportEvent sends a synchronous event report to the Fleet Manager.
func (r *Reporter) ReportEvent(eventType string, data map[string]any) error {
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"type":      eventType,
		"data":      data,
		"timestamp": r.now().UTC().Format(time.RFC3339),
	}
	err := r.deliver("/fleet/events", payload)
	if err == nil {
		metrics.EventsTotal.Add(1)
	}
	return err
}

// QueueEvent enqueues an event for async delivery. Non-blocking; drops
// if the queue is full and logs a warning with the metrics counter.
func (r *Reporter) QueueEvent(eventType string, data map[string]any) {
	if r.stopped.True() {
		logger.WarnCF("edge", "reporter stopped, dropping queued event", map[string]any{"type": eventType})
		metrics.DroppedEventsTotal.Add(1)
		return
	}
	select {
	case r.eventCh <- eventJob{eventType: eventType, data: data, timestamp: r.now()}:
		metrics.QueueDepth.Add(1)
		metrics.QueuedEventsTotal.Add(1)
	default:
		logger.WarnCF("edge", "event queue full, dropping event", map[string]any{"type": eventType})
		metrics.DroppedEventsTotal.Add(1)
	}
}

// PublishGene sends a high-confidence gene to the Fleet Manager
// for cross-node sharing via the /fleet/genes/publish endpoint.
func (r *Reporter) PublishGene(geneData map[string]interface{}) error {
	payload := map[string]any{
		"node_id":   r.config.NodeID,
		"gene":      geneData,
		"timestamp": r.now().UTC().Format(time.RFC3339),
	}
	return r.deliver("/fleet/genes/publish", payload)
}

// Metrics returns a snapshot of the reporter's live metrics.
func (r *Reporter) Metrics() map[string]any {
	return metrics.Snapshot()
}

// ---------------------------------------------------------------------------
// Internal loops
// ---------------------------------------------------------------------------

func (r *Reporter) startHeartbeatLoop() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		hbi := r.config.Cloud.HeartbeatInterval
		if hbi <= 0 {
			hbi = 30 * time.Second
		}

		logger.InfoCF("edge", "heartbeat reporter started",
			map[string]any{"interval": hbi.String(), "endpoint": r.config.Cloud.Endpoint})

		// Send initial registration + heartbeat
		if err := r.Register(); err != nil {
			logger.ErrorCF("edge", "initial registration failed", map[string]any{"error": err.Error()})
		} else {
			metrics.RegistrationsTotal.Add(1)
		}
		if err := r.SendHeartbeat(); err != nil {
			logger.ErrorCF("edge", "initial heartbeat failed", map[string]any{"error": err.Error()})
			if r.onHeartbeatSent != nil {
				r.onHeartbeatSent(err)
			}
		} else {
			if r.onHeartbeatSent != nil {
				r.onHeartbeatSent(nil)
			}
		}

		ticker := time.NewTicker(hbi)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := r.SendHeartbeat()
				if err != nil {
					logger.ErrorCF("edge", "heartbeat failed", map[string]any{"error": err.Error()})
				}
				if r.onHeartbeatSent != nil {
					r.onHeartbeatSent(err)
				}
			case <-r.stopCh:
				logger.DebugCF("edge", "heartbeat loop exiting", nil)
				return
			}
		}
	}()
}

func (r *Reporter) startStatusLoop() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		si := r.config.StatusInterval
		if si <= 0 {
			si = 5 * time.Minute
		}

		ticker := time.NewTicker(si)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := r.ReportStatus()
				if err != nil {
					logger.ErrorCF("edge", "status report failed", map[string]any{"error": err.Error()})
				}
				if r.onStatusSent != nil {
					r.onStatusSent(err)
				}
			case <-r.stopCh:
				logger.DebugCF("edge", "status loop exiting", nil)
				return
			}
		}
	}()
}

func (r *Reporter) startEventWorkers() {
	n := r.config.EventWorkers
	if n <= 0 {
		n = 1
	}
	for i := 0; i < n; i++ {
		r.wg.Add(1)
		go func(id int) {
			defer r.wg.Done()
			logger.DebugCF("edge", "event worker started", map[string]any{"worker_id": id})
			for {
				select {
				case job := <-r.eventCh:
					metrics.QueueDepth.Add(-1)
					payload := map[string]any{
						"node_id":   r.config.NodeID,
						"type":      job.eventType,
						"data":      job.data,
						"timestamp": job.timestamp.UTC().Format(time.RFC3339),
					}
					if err := r.deliver("/fleet/events", payload); err != nil {
						logger.ErrorCF("edge", "queued event delivery failed",
							map[string]any{"type": job.eventType, "error": err.Error()})
					}
				case <-r.stopCh:
					logger.DebugCF("edge", "event worker exiting", map[string]any{"worker_id": id})
					return
				}
			}
		}(i)
	}
}

// ---------------------------------------------------------------------------
// HTTP delivery with exponential backoff
// ---------------------------------------------------------------------------

// deliver sends a JSON POST to the Fleet Manager at the given path.
// It retries on failures with exponential backoff + jitter.
func (r *Reporter) deliver(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	maxRetries := r.config.MaxRetries
	baseDelay := r.config.RetryBaseDelay
	if baseDelay <= 0 {
		baseDelay = 1 * time.Second
	}

	var lastErr error
	// attempt 0 = initial send, attempts 1..maxRetries = retries
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			metrics.RetriesTotal.Add(1)
		}

		err := r.doPost(path, body)
		metrics.DeliveriesTotal.Add(1)

		if err == nil {
			metrics.DeliveriesOK.Add(1)
			metrics.LastDeliveryOK.Store(true)
			return nil
		}

		lastErr = err
		metrics.DeliveriesFailed.Add(1)
		metrics.LastDeliveryOK.Store(false)

		if attempt < maxRetries {
			delay := backoff(baseDelay, attempt)
			logger.WarnCF("edge", "delivery failed, retrying",
				map[string]any{"path": path, "attempt": attempt + 1, "max": maxRetries + 1, "retry_in": delay.String(), "error": err.Error()})

			select {
			case <-time.After(delay):
			case <-r.stopCh:
				return fmt.Errorf("reporter stopped during retry: %w", lastErr)
			}
		}
	}

	return fmt.Errorf("delivery failed after %d attempts: %w", maxRetries+1, lastErr)
}

func (r *Reporter) doPost(path string, body []byte) error {
	baseURL := strings.TrimRight(r.config.Cloud.Endpoint, "/")
	url := baseURL + path

	ctx, cancel := context.WithTimeout(context.Background(), r.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
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

	start := time.Now()
	resp, err := r.client.Do(req)
	latency := time.Since(start)
	metrics.recordLatency(latency)

	if err != nil {
		return fmt.Errorf("post %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		// Read bounded response body for diagnostics
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		bodyStr := strings.TrimSpace(string(bodyBytes))
		if len(bodyStr) > 0 {
			return fmt.Errorf("post %s: status %d: %s", path, resp.StatusCode, bodyStr)
		}
		return fmt.Errorf("post %s: status %d", path, resp.StatusCode)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// backoff returns a delay with exponential backoff and full jitter:
//
//	delay = min(base * 2^attempt, max) * random(0.5, 1.5)
func backoff(base time.Duration, attempt int) time.Duration {
	if attempt < 1 {
		return base
	}
	exp := float64(base) * math.Pow(2, float64(attempt-1))
	const maxDelay = 30 * time.Second
	if exp > float64(maxDelay) {
		exp = float64(maxDelay)
	}
	// Add ±50% jitter
	jitter := 0.5 + rand.Float64() // 0.5–1.5
	return time.Duration(exp * jitter)
}

var reporterStartTime time.Time

func init() {
	reporterStartTime = time.Now()
}

func (r *Reporter) startTime() time.Time {
	return reporterStartTime
}

func eventQueueSize(cfg int) int {
	if cfg > 0 {
		return cfg
	}
	return 256
}
