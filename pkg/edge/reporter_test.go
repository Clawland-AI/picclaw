package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReporterPostsStatusAndEventsWithAuth(t *testing.T) {
	type receivedRequest struct {
		Path string
		Auth string
		Body map[string]any
	}

	requests := make([]receivedRequest, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		requests = append(requests, receivedRequest{
			Path: r.URL.Path,
			Auth: r.Header.Get("Authorization"),
			Body: body,
		})
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	reporter := NewReporterWithClient(Config{
		NodeID:   "rack-a1",
		NodeName: "Rack A1",
		Cloud: CloudConfig{
			Endpoint:          server.URL,
			Token:             "fleet-token",
			HeartbeatInterval: time.Minute,
		},
	}, server.Client())

	if err := reporter.ReportStatus("degraded", map[string]any{"temperature_c": 41.5}); err != nil {
		t.Fatalf("report status: %v", err)
	}
	if err := reporter.ReportEvent("threshold_crossed", map[string]any{"sensor": "temperature"}); err != nil {
		t.Fatalf("report event: %v", err)
	}

	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}
	if requests[0].Path != "/fleet/status" {
		t.Fatalf("expected status path, got %q", requests[0].Path)
	}
	if requests[1].Path != "/fleet/events" {
		t.Fatalf("expected event path, got %q", requests[1].Path)
	}
	for _, req := range requests {
		if req.Auth != "Bearer fleet-token" {
			t.Fatalf("expected bearer auth header, got %q", req.Auth)
		}
		if req.Body["node_id"] != "rack-a1" {
			t.Fatalf("expected node_id rack-a1, got %#v", req.Body["node_id"])
		}
		if _, ok := req.Body["timestamp"].(string); !ok {
			t.Fatalf("expected timestamp in request body, got %#v", req.Body)
		}
	}
	if requests[0].Body["status"] != "degraded" {
		t.Fatalf("expected status degraded, got %#v", requests[0].Body["status"])
	}
	metrics, ok := requests[0].Body["metrics"].(map[string]any)
	if !ok || metrics["temperature_c"] != 41.5 {
		t.Fatalf("expected status metrics, got %#v", requests[0].Body["metrics"])
	}
	if requests[1].Body["type"] != "threshold_crossed" {
		t.Fatalf("expected event type threshold_crossed, got %#v", requests[1].Body["type"])
	}
}
