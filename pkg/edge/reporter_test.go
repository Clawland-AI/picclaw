package edge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type fleetHandler struct {
	mu           sync.Mutex
	requests     []*http.Request
	bodies       []map[string]any
	statusCode   int
	responseBody string
}

func newFleetHandler() *fleetHandler {
	return &fleetHandler{statusCode: 200, responseBody: `{"status":"ok"}`}
}

func (fh *fleetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.mu.Lock()
	defer fh.mu.Unlock()
	var body map[string]any
	json.NewDecoder(r.Body).Decode(&body)
	fh.requests = append(fh.requests, r)
	fh.bodies = append(fh.bodies, body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fh.statusCode)
	fmt.Fprint(w, fh.responseBody)
}
func (fh *fleetHandler) last() (*http.Request, map[string]any) {
	fh.mu.Lock()
	defer fh.mu.Unlock()
	if len(fh.requests) == 0 {
		return nil, nil
	}
	return fh.requests[len(fh.requests)-1], fh.bodies[len(fh.bodies)-1]
}
func newTestReporter(fh *fleetHandler) (*Reporter, *httptest.Server) {
	ts := httptest.NewServer(fh)
	c := DefaultConfig()
	c.Cloud.Endpoint = ts.URL
	c.NodeID = "test-node-001"
	c.NodeName = "test-edge"
	c.Cloud.Token = "test-token-abc"
	return NewReporter(c), ts
}

func TestRegister(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()
	if err := r.Register(); err != nil {
		t.Fatal(err)
	}
	lr, body := fh.last()
	if lr == nil || lr.URL.Path != "/fleet/register" {
		t.Fatalf("expected /fleet/register, got %v", lr)
	}
	if body["node_id"] != "test-node-001" || body["type"] != "picclaw" {
		t.Fatalf("unexpected body: %v", body)
	}
	if lr.Header.Get("Authorization") != "Bearer test-token-abc" {
		t.Errorf("bad auth: %q", lr.Header.Get("Authorization"))
	}
}
func TestSendHeartbeat(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()
	if err := r.SendHeartbeat(); err != nil {
		t.Fatal(err)
	}
	if _, body := fh.last(); body["status"] != "online" {
		t.Errorf("expected online, got %v", body["status"])
	}
}
func TestReportStatus(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()
	if err := r.ReportStatus(); err != nil {
		t.Fatal(err)
	}
	_, body := fh.last()
	if body["status"] != "online" || body["node_name"] != "test-edge" {
		t.Errorf("unexpected body: %v", body)
	}
}
func TestReportEvent(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()
	if err := r.ReportEvent("sensor_reading", map[string]any{"temperature": 22.5}); err != nil {
		t.Fatal(err)
	}
	_, body := fh.last()
	if body["type"] != "sensor_reading" || body["data"] == nil {
		t.Errorf("unexpected body: %v", body)
	}
}
func TestNoAuthWhenTokenEmpty(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Cloud.Token = ""
	cfg.NodeID = "no-token-node"
	fh := newFleetHandler()
	ts := httptest.NewServer(fh)
	defer ts.Close()
	cfg.Cloud.Endpoint = ts.URL
	r := NewReporter(cfg)
	r.SendHeartbeat()
	lr, _ := fh.last()
	if lr == nil {
		t.Fatal("expected request")
	}
	if auth := lr.Header.Get("Authorization"); auth != "" {
		t.Errorf("expected no auth, got %q", auth)
	}
	if lr.Header.Get("X-Clawland-Node-ID") != "no-token-node" {
		t.Errorf("bad node id header")
	}
}
func TestHTTPFailure(t *testing.T) {
	fh := newFleetHandler()
	fh.statusCode = 500
	fh.responseBody = `{"error":"internal error"}`
	r, ts := newTestReporter(fh)
	defer ts.Close()
	if err := r.SendHeartbeat(); err == nil || !contains(err.Error(), "500") {
		t.Fatalf("expected 500 error, got: %v", err)
	}
}
func TestIdempotentStop(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()
	r.Start()
	time.Sleep(100 * time.Millisecond)
	r.Stop()
	r.Stop()
	r.Stop()
}
func TestHeartbeatWithGeneStats(t *testing.T) {
	fh := newFleetHandler()
	r, ts := newTestReporter(fh)
	defer ts.Close()
	r.SetGeneProvider(&mockGeneProvider{})
	if err := r.SendHeartbeat(); err != nil {
		t.Fatal(err)
	}
	_, body := fh.last()
	if _, ok := body["gene_stats"]; !ok {
		t.Error("expected gene_stats in heartbeat")
	}
}

type mockGeneProvider struct{}
func (m *mockGeneProvider) GetStats() map[string]interface{} {
	return map[string]interface{}{"total_genes": 10, "verified_genes": 5, "avg_confidence": 0.85}
}
func (m *mockGeneProvider) GetHighConfidenceGenes(float64, int) []map[string]interface{} { return nil }

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
