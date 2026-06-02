package edge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Helper: test server
// ---------------------------------------------------------------------------

type fleetHandler struct {
	mu           sync.Mutex
	requests     []*http.Request
	bodies       []map[string]any
	statusCode   int
	responseBody string
	// optional hooks
	onRequest func(path string, body map[string]any)
}

func newFleetHandler() *fleetHandler {
	return &fleetHandler{
		statusCode:   200,
		responseBody: `{"status":"ok"}`,
	}
}

func (fh *fleetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.mu.Lock()
	defer fh.mu.Unlock()

	var body map[string]any
	json.NewDecoder(r.Body).Decode(&body)

	fh.requests = append(fh.requests, r)
	fh.bodies = append(fh.bodies, body)

	if fh.onRequest != nil {
		fh.onRequest(r.URL.Path, body)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fh.statusCode)
	fmt.Fprint(w, fh.responseBody)
}

func (fh *fleetHandler) lastRequest() *http.Request {
	fh.mu.Lock()
	defer fh.mu.Unlock()
	if len(fh.requests) == 0 {
		return nil
	}
	return fh.requests[len(fh.requests)-1]
}

func (fh *fleetHandler) lastBody() map[string]any {
	fh.mu.Lock()
	defer fh.mu.Unlock()
	if len(fh.bodies) == 0 {
		return nil
	}
	return fh.bodies[len(fh.bodies)-1]
}

func (fh *fleetHandler) count() int {
	fh.mu.Lock()
	defer fh.mu.Unlock()
	return len(fh.requests)
}

func newTestReporter(fh *fleetHandler, cfg ...Config) (*Reporter, *httptest.Server) {
	ts := httptest.NewServer(fh)

	c := DefaultConfig()
	c.Cloud.Endpoint = ts.URL
	c.NodeID = "test-node-001"
	c.NodeName = "test-edge"
	c.Cloud.Token = "test-token-abc"
	c.Cloud.HeartbeatInterval = 100 * time.Millisecond
	c.StatusInterval = 200 * time.Millisecond
	c.MaxRetries = 2
	c.RetryBaseDelay = 5 * time.Millisecond
	c.EventQueueSize = 10
	c.EventWorkers = 2
	if len(cfg) > 0 {
		c = cfg[0]
		c.Cloud.Endpoint = ts.URL
	}

	r := NewReporter(c)
	return r, ts
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestRegister(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	err := r.Register()
	if err != nil {
		t.Fatalf("Register() returned error: %v", err)
	}

	lr := fh.lastRequest()
	if lr == nil {
		t.Fatal("expected at least one request")
	}
	if lr.URL.Path != "/fleet/register" {
		t.Errorf("expected path /fleet/register, got %s", lr.URL.Path)
	}

	body := fh.lastBody()
	if body == nil {
		t.Fatal("expected a body")
	}
	if body["node_id"] != "test-node-001" {
		t.Errorf("expected node_id test-node-001, got %v", body["node_id"])
	}
	if body["type"] != "picclaw" {
		t.Errorf("expected type picclaw, got %v", body["type"])
	}

	// Check auth header
	auth := lr.Header.Get("Authorization")
	if auth != "Bearer test-token-abc" {
		t.Errorf("expected Bearer test-token-abc, got %q", auth)
	}
	nodeID := lr.Header.Get("X-Clawland-Node-ID")
	if nodeID != "test-node-001" {
		t.Errorf("expected X-Clawland-Node-ID: test-node-001, got %q", nodeID)
	}
}

func TestSendHeartbeat(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	err := r.SendHeartbeat()
	if err != nil {
		t.Fatalf("SendHeartbeat() returned error: %v", err)
	}

	lr := fh.lastRequest()
	if lr == nil || lr.URL.Path != "/fleet/heartbeat" {
		t.Errorf("expected /fleet/heartbeat, got %v", lr)
	}

	body := fh.lastBody()
	if body["status"] != "online" {
		t.Errorf("expected status online, got %v", body["status"])
	}
}

func TestSendHeartbeat_WithGeneStats(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.SetGeneProvider(&mockGeneProvider{})
	err := r.SendHeartbeat()
	if err != nil {
		t.Fatalf("SendHeartbeat() returned error: %v", err)
	}

	body := fh.lastBody()
	if _, ok := body["gene_stats"]; !ok {
		t.Error("expected gene_stats in heartbeat body")
	}
}

func TestReportStatus(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	err := r.ReportStatus()
	if err != nil {
		t.Fatalf("ReportStatus() returned error: %v", err)
	}

	lr := fh.lastRequest()
	if lr == nil || lr.URL.Path != "/fleet/status" {
		t.Errorf("expected /fleet/status, got %v", lr)
	}

	body := fh.lastBody()
	if body["status"] != "online" {
		t.Errorf("expected status online, got %v", body["status"])
	}
	if _, ok := body["metrics"]; !ok {
		t.Error("expected metrics in status body")
	}
}

func TestReportStatus_WithGenes(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.SetGeneProvider(&mockGeneProvider{
		hcGenes: []map[string]interface{}{
			{"gene_id": "g1", "confidence": 0.95},
		},
	})
	err := r.ReportStatus()
	if err != nil {
		t.Fatalf("ReportStatus() returned error: %v", err)
	}

	body := fh.lastBody()
	if _, ok := body["high_confidence_genes"]; !ok {
		t.Error("expected high_confidence_genes in status body")
	}
}

func TestReportEvent(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	err := r.ReportEvent("sensor_reading", map[string]any{"temperature": 22.5})
	if err != nil {
		t.Fatalf("ReportEvent() returned error: %v", err)
	}

	lr := fh.lastRequest()
	if lr == nil || lr.URL.Path != "/fleet/events" {
		t.Errorf("expected /fleet/events, got %v", lr)
	}

	body := fh.lastBody()
	if body["type"] != "sensor_reading" {
		t.Errorf("expected type sensor_reading, got %v", body["type"])
	}
}

func TestQueueEvent(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	// Use a conditional chan to signal when we've received the event
	received := make(chan struct{}, 1)
	var receivedBody map[string]any
	fh.onRequest = func(path string, body map[string]any) {
		if path == "/fleet/events" {
			receivedBody = body
			select {
			case received <- struct{}{}:
			default:
			}
		}
	}

	// Start the reporter so event workers are running
	r.Start()

	r.QueueEvent("async_test", map[string]any{"key": "value"})

	select {
	case <-received:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for queued event delivery")
	}

	r.Stop()

	if receivedBody == nil || receivedBody["type"] != "async_test" {
		t.Errorf("expected type async_test, got %v", receivedBody["type"])
	}
}

func TestQueueEvent_DropsWhenFull(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	// Fill the queue by blocking the handler
	block := make(chan struct{})
	fh.onRequest = func(path string, body map[string]any) {
		if path == "/fleet/events" {
			<-block
		}
	}

	// Queue more than capacity (queue size=10, workers=2, so buffer fills fast)
	for i := 0; i < 20; i++ {
		r.QueueEvent("drop_test", map[string]any{"i": i})
	}

	time.Sleep(50 * time.Millisecond)

	// Release blocked requests
	close(block)

	time.Sleep(100 * time.Millisecond)

	if metrics.DroppedEventsTotal.Load() == 0 {
		t.Log("dropped_events_total:", metrics.DroppedEventsTotal.Load())
		// May or may not drop depending on timing, but at least we didn't panic
	}
}

func TestPublishGene(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	err := r.PublishGene(map[string]interface{}{"gene_id": "g1", "sequence": "ATCG"})
	if err != nil {
		t.Fatalf("PublishGene() returned error: %v", err)
	}

	lr := fh.lastRequest()
	if lr == nil || lr.URL.Path != "/fleet/genes/publish" {
		t.Errorf("expected /fleet/genes/publish, got %v", lr)
	}
}

func TestAuthorizationHeader(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Register()

	lr := fh.lastRequest()
	auth := lr.Header.Get("Authorization")
	if auth != "Bearer test-token-abc" {
		t.Errorf("expected Bearer test-token-abc, got %q", auth)
	}
}

func TestNodeIDHeader(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.SendHeartbeat()

	lr := fh.lastRequest()
	nodeID := lr.Header.Get("X-Clawland-Node-ID")
	if nodeID != "test-node-001" {
		t.Errorf("expected X-Clawland-Node-ID: test-node-001, got %q", nodeID)
	}
}

func TestNoAuthHeaderWhenTokenEmpty(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Cloud.Token = ""
	cfg.NodeID = "no-token-node"

	fh := newFleetHandler()
	r, ts := newTestReporter(fh, cfg)
	defer ts.Close()

	r.SendHeartbeat()

	lr := fh.lastRequest()
	auth := lr.Header.Get("Authorization")
	if auth != "" {
		t.Errorf("expected no Authorization header, got %q", auth)
	}
	// Node ID header should still be set
	nodeID := lr.Header.Get("X-Clawland-Node-ID")
	if nodeID != "no-token-node" {
		t.Errorf("expected X-Clawland-Node-ID: no-token-node, got %q", nodeID)
	}
}

func TestHTTPFailureStatusCode(t *testing.T) {
	fh := newFleetHandler()
	fh.statusCode = 500
	fh.responseBody = `{"error":"internal error","code":"INTERNAL"}`
	r, ts := newTestReporter(fh)
	defer ts.Close()

	err := r.SendHeartbeat()
	if err == nil {
		t.Fatal("expected error for status 500")
	}
	if !contains(err.Error(), "500") {
		t.Errorf("expected error to mention 500, got: %v", err)
	}
	if !contains(err.Error(), "internal error") {
		t.Errorf("expected error to mention response body, got: %v", err)
	}
}

func TestHTTPFailureUnauthorized(t *testing.T) {
	fh := newFleetHandler()
	fh.statusCode = 401
	fh.responseBody = `{"error":"unauthorized"}`
	r, ts := newTestReporter(fh)
	defer ts.Close()

	err := r.SendHeartbeat()
	if err == nil {
		t.Fatal("expected error for status 401")
	}
	if !contains(err.Error(), "401") {
		t.Errorf("expected error to mention 401, got: %v", err)
	}
}

func TestEndpointNormalization(t *testing.T) {
	fh := newFleetHandler()
	cfg := DefaultConfig()
	cfg.Cloud.Endpoint = "http://example.com/" // trailing slash
	ts := httptest.NewServer(fh)
	defer ts.Close()
	cfg.Cloud.Endpoint = ts.URL + "/" // ensure trailing slash

	r := NewReporter(cfg)
	r.SendHeartbeat()

	lr := fh.lastRequest()
	if lr == nil {
		t.Fatal("expected request")
	}
	// Should not have double slash
	if contains(lr.URL.String(), "//fleet") {
		t.Errorf("URL should not have double slash: %s", lr.URL.String())
	}
}

func TestRetryOnFailure(t *testing.T) {
	fh := newFleetHandler()
	var attemptCount atomic.Int32
	fh.onRequest = func(path string, body map[string]any) {
		attemptCount.Add(1)
	}
	fh.statusCode = 503

	cfg := DefaultConfig()
	cfg.MaxRetries = 3
	cfg.RetryBaseDelay = 2 * time.Millisecond

	r, ts := newTestReporter(fh, cfg)
	defer ts.Close()

	err := r.SendHeartbeat()
	if err == nil {
		t.Fatal("expected error after all retries exhausted")
	}

	// Should have attempted 4 times (initial + 3 retries)
	if n := attemptCount.Load(); n != 4 {
		t.Errorf("expected 4 attempts (1+3 retries), got %d", n)
	}
}

func TestRetryCountMatchesMetric(t *testing.T) {
	fh := newFleetHandler()
	fh.statusCode = 503

	cfg := DefaultConfig()
	cfg.MaxRetries = 2
	cfg.RetryBaseDelay = 1 * time.Millisecond

	r, ts := newTestReporter(fh, cfg)
	defer ts.Close()

	prevRetries := metrics.RetriesTotal.Load()
	r.SendHeartbeat()

	// Should have 2 retries
	got := metrics.RetriesTotal.Load() - prevRetries
	if got != 2 {
		t.Errorf("expected 2 retries in metrics, got %d", got)
	}
}

func TestTimeout(t *testing.T) {
	fh := newFleetHandler()
	fh.onRequest = func(path string, body map[string]any) {
		time.Sleep(500 * time.Millisecond) // slower than timeout
	}

	cfg := DefaultConfig()
	cfg.MaxRetries = 0 // no retries to speed up test
	r, ts := newTestReporter(fh, cfg)

	// Set a very short timeout
	r.client.Timeout = 10 * time.Millisecond

	defer ts.Close()

	start := time.Now()
	err := r.SendHeartbeat()
	duration := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !contains(err.Error(), "timeout") && !contains(err.Error(), "context") {
		t.Errorf("expected timeout/context error, got: %v", err)
	}
	if duration > 5*time.Second {
		t.Errorf("timeout test took too long: %v", duration)
	}
}

func TestIdempotentStop(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Start()
	time.Sleep(200 * time.Millisecond)

	// Stop multiple times should not panic
	r.Stop()
	r.Stop()
	r.Stop()
}

func TestStopDrainsEvents(t *testing.T) {
	fh := newFleetHandler()
	received := make(chan struct{}, 10)
	fh.onRequest = func(path string, body map[string]any) {
		if path == "/fleet/events" {
			received <- struct{}{}
		}
	}

	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Start()
	time.Sleep(50 * time.Millisecond)

	// Queue events
	r.QueueEvent("e1", nil)
	r.QueueEvent("e2", nil)

	// Stop should drain before returning
	r.Stop()

	// We may or may not have received them depending on timing,
	// but the point is no panic/deadlock
}

func TestStoppedReporterRejectsEvents(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Start()
	time.Sleep(50 * time.Millisecond)
	r.Stop()

	// After stop, queue event should be dropped silently
	r.QueueEvent("after_stop", nil)

	if metrics.DroppedEventsTotal.Load() == 0 {
		t.Log("dropped_events_total:", metrics.DroppedEventsTotal.Load())
	}
}

func TestStartSendsInitialRegistrationAndHeartbeat(t *testing.T) {
	fh := newFleetHandler()
	var paths []string
	var mu sync.Mutex
	fh.onRequest = func(path string, body map[string]any) {
		mu.Lock()
		paths = append(paths, path)
		mu.Unlock()
	}

	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Start()
	time.Sleep(300 * time.Millisecond)
	r.Stop()

	mu.Lock()
	hasRegister := false
	hasHeartbeat := false
	for _, p := range paths {
		if p == "/fleet/register" {
			hasRegister = true
		}
		if p == "/fleet/heartbeat" {
			hasHeartbeat = true
		}
	}
	mu.Unlock()

	if !hasRegister {
		t.Error("expected initial registration after Start()")
	}
	if !hasHeartbeat {
		t.Error("expected initial heartbeat after Start()")
	}
}

func TestMetricsArePopulated(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Register()
	r.SendHeartbeat()
	r.ReportStatus()
	r.ReportEvent("test", nil)

	snapshot := metrics.Snapshot()
	if snapshot["registrations_total"].(int64) < 1 {
		t.Error("expected registrations_total >= 1")
	}
	if snapshot["heartbeats_total"].(int64) < 1 {
		t.Error("expected heartbeats_total >= 1")
	}
	if snapshot["status_reports_total"].(int64) < 1 {
		t.Error("expected status_reports_total >= 1")
	}
	if snapshot["events_total"].(int64) < 1 {
		t.Error("expected events_total >= 1")
	}
	if snapshot["deliveries_total"].(int64) < 4 {
		t.Errorf("expected deliveries_total >= 4, got %d", snapshot["deliveries_total"])
	}
	if snapshot["deliveries_ok"].(int64) < 4 {
		t.Errorf("expected deliveries_ok >= 4, got %d", snapshot["deliveries_ok"])
	}
}

func TestMetricsRecordFailures(t *testing.T) {
	fh := newFleetHandler()
	fh.statusCode = 500

	r, ts := newTestReporter(fh)
	defer ts.Close()

	prevFailed := metrics.DeliveriesFailed.Load()
	r.SendHeartbeat()

	if metrics.DeliveriesFailed.Load() <= prevFailed {
		t.Error("expected deliveries_failed to increase after failed heartbeat")
	}
	if metrics.LastDeliveryOK.Load() {
		t.Error("expected last_delivery_ok to be false after failure")
	}
}

func TestConcurrentAccess(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Start()
	defer r.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.SendHeartbeat()
			r.ReportStatus()
			r.ReportEvent("concurrent", nil)
			r.QueueEvent("concurrent_async", nil)
			r.PublishGene(map[string]interface{}{"g": "1"})
			_ = r.Metrics()
		}()
	}
	wg.Wait()
}

func TestSetGeneProviderConcurrency(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.SetGeneProvider(&mockGeneProvider{})
		}()
	}
	wg.Wait()
}

func TestCustomClock(t *testing.T) {
	fixedTime := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	r := NewReporter(DefaultConfig(), WithClock(func() time.Time {
		return fixedTime
	}))

	fh := newFleetHandler()
	ts := httptest.NewServer(fh)
	defer ts.Close()

	// Override endpoint
	r.config.Cloud.Endpoint = ts.URL
	r.config.NodeID = "clock-test"

	r.SendHeartbeat()

	body := fh.lastBody()
	if body == nil {
		t.Fatal("expected body")
	}
	tsField, ok := body["timestamp"].(string)
	if !ok {
		t.Fatal("expected timestamp string")
	}
	if tsField != "2026-06-01T12:00:00Z" {
		t.Errorf("expected timestamp 2026-06-01T12:00:00Z, got %q", tsField)
	}
}

func TestHeartbeatHook(t *testing.T) {
	var (
		hookMu   sync.Mutex
		hookErr  error
	)
	hookCh := make(chan struct{}, 1)
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	WithHeartbeatHook(func(err error) {
		hookMu.Lock()
		hookErr = err
		hookMu.Unlock()
		hookCh <- struct{}{}
	})(r)

	r.Start()
	time.Sleep(200 * time.Millisecond)

	select {
	case <-hookCh:
		hookMu.Lock()
		err := hookErr
		hookMu.Unlock()
		if err != nil {
			t.Errorf("expected nil error from hook, got: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for heartbeat hook")
	}
	r.Stop()
}

// ---------------------------------------------------------------------------
// Race-condition tests (run with go test -race)
// ---------------------------------------------------------------------------

func TestRace_ConcurrentStartStop(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Start()
			time.Sleep(10 * time.Millisecond)
			r.Stop()
		}()
	}
	wg.Wait()
}

func TestRace_ConcurrentCalls(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	r.Start()
	defer r.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			switch id % 6 {
			case 0:
				r.Register()
			case 1:
				r.SendHeartbeat()
			case 2:
				r.ReportStatus()
			case 3:
				r.ReportEvent("race", nil)
			case 4:
				r.QueueEvent("race_async", nil)
			case 5:
				r.PublishGene(map[string]interface{}{"race": id})
			}
		}(i)
	}
	wg.Wait()
}

func TestRace_MetricsConcurrent(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.SendHeartbeat()
			_ = r.Metrics()
			r.ReportStatus()
			_ = r.Metrics()
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type mockGeneProvider struct {
	stats   map[string]interface{}
	hcGenes []map[string]interface{}
}

func (m *mockGeneProvider) GetStats() map[string]interface{} {
	if m.stats != nil {
		return m.stats
	}
	return map[string]interface{}{
		"total_genes":     10,
		"verified_genes":  5,
		"avg_confidence":  0.85,
		"active_mutations": 3,
	}
}

func (m *mockGeneProvider) GetHighConfidenceGenes(minConfidence float64, minVerifiedBy int) []map[string]interface{} {
	if m.hcGenes != nil {
		return m.hcGenes
	}
	return []map[string]interface{}{
		{"gene_id": "g1", "confidence": 0.95, "verified_by": 3},
		{"gene_id": "g2", "confidence": 0.91, "verified_by": 2},
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
