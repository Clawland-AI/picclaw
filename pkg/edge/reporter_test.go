package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type capturedReporterRequest struct {
	path          string
	authorization string
	nodeHeader    string
	body          map[string]any
}

type fakeGeneProvider struct{}

func (fakeGeneProvider) GetStats() map[string]interface{} {
	return map[string]interface{}{"genes": 2, "capsules": 5}
}

func (fakeGeneProvider) GetHighConfidenceGenes(float64, int) []map[string]interface{} {
	return []map[string]interface{}{{"id": "gene-1", "confidence": 0.92}}
}

func TestReporterPostsHeartbeatStatusAndEvents(t *testing.T) {
	var (
		mu       sync.Mutex
		requests []capturedReporterRequest
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		requests = append(requests, capturedReporterRequest{
			path:          r.URL.Path,
			authorization: r.Header.Get("Authorization"),
			nodeHeader:    r.Header.Get("X-Clawland-Node-ID"),
			body:          body,
		})
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	reporter := NewReporter(Config{
		NodeID:   "edge-01",
		NodeName: "Rack A1",
		Cloud: CloudConfig{
			Endpoint: server.URL + "/",
			Token:    "secret-token",
		},
	})
	reporter.SetGeneProvider(fakeGeneProvider{})

	reporter.sendHeartbeat()
	if err := reporter.ReportStatus("degraded", map[string]any{"reason": "temperature"}); err != nil {
		t.Fatalf("ReportStatus failed: %v", err)
	}
	if err := reporter.ReportEvent("sensor.alert", map[string]any{"temp_c": 31.5}); err != nil {
		t.Fatalf("ReportEvent failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(requests) != 3 {
		t.Fatalf("expected 3 reports, got %d", len(requests))
	}
	for _, req := range requests {
		if req.authorization != "Bearer secret-token" {
			t.Fatalf("expected bearer token on %s, got %q", req.path, req.authorization)
		}
		if req.nodeHeader != "edge-01" {
			t.Fatalf("expected node header on %s, got %q", req.path, req.nodeHeader)
		}
		if req.body["node_id"] != "edge-01" {
			t.Fatalf("expected node_id in %s payload, got %#v", req.path, req.body["node_id"])
		}
	}

	heartbeat := findReporterRequest(t, requests, "/fleet/heartbeat")
	if heartbeat.body["gene_stats"] == nil {
		t.Fatalf("expected heartbeat to include gene stats: %#v", heartbeat.body)
	}

	status := findReporterRequest(t, requests, "/fleet/status")
	if status.body["status"] != "degraded" {
		t.Fatalf("expected degraded status, got %#v", status.body["status"])
	}
	if status.body["high_confidence_genes"] == nil {
		t.Fatalf("expected status to include high-confidence genes: %#v", status.body)
	}

	event := findReporterRequest(t, requests, "/fleet/events")
	if event.body["type"] != "sensor.alert" {
		t.Fatalf("expected sensor alert event, got %#v", event.body["type"])
	}
}

func TestReporterQueueEventFlushesFromLoop(t *testing.T) {
	paths := make(chan string, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths <- r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	reporter := NewReporter(Config{
		NodeID:   "edge-02",
		NodeName: "Rack A2",
		Cloud: CloudConfig{
			Endpoint:          server.URL,
			HeartbeatInterval: time.Hour,
			StatusInterval:    time.Hour,
		},
	})
	defer reporter.Stop()

	if !reporter.QueueEvent("command.completed", map[string]any{"command_id": "cmd-1"}) {
		t.Fatal("expected event to be queued")
	}
	go reporter.StartHeartbeat()

	timeout := time.After(time.Second)
	for {
		select {
		case path := <-paths:
			if path == "/fleet/events" {
				return
			}
		case <-timeout:
			t.Fatal("timed out waiting for queued event report")
		}
	}
}

func TestNewReporterAppliesDefaultIntervals(t *testing.T) {
	reporter := NewReporter(Config{})
	if reporter.config.Cloud.HeartbeatInterval != defaultHeartbeatInterval {
		t.Fatalf("expected default heartbeat interval, got %s", reporter.config.Cloud.HeartbeatInterval)
	}
	if reporter.config.Cloud.StatusInterval != defaultStatusInterval {
		t.Fatalf("expected default status interval, got %s", reporter.config.Cloud.StatusInterval)
	}
}

func findReporterRequest(t *testing.T, requests []capturedReporterRequest, path string) capturedReporterRequest {
	t.Helper()
	for _, req := range requests {
		if req.path == path {
			return req
		}
	}
	t.Fatalf("missing request for %s in %#v", path, requests)
	return capturedReporterRequest{}
}
