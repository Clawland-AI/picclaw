package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type reporterRequest struct {
	Path          string
	Authorization string
	Payload       map[string]any
}

type fakeGeneStats struct{}

func (fakeGeneStats) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"genes": float64(3),
	}
}

func (fakeGeneStats) GetHighConfidenceGenes(minConfidence float64, minVerifiedBy int) []map[string]interface{} {
	return []map[string]interface{}{}
}

func TestReporterSendsRegistrationHeartbeatStatusAndEvents(t *testing.T) {
	var requests []reporterRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		requests = append(requests, reporterRequest{
			Path:          r.URL.Path,
			Authorization: r.Header.Get("Authorization"),
			Payload:       payload,
		})
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	reporter := NewReporter(Config{
		NodeID:   "edge-01",
		NodeName: "Pond Edge 01",
		Cloud: CloudConfig{
			Endpoint:          server.URL + "/",
			HeartbeatInterval: time.Second,
			Token:             "fleet-token",
		},
	})
	reporter.SetGeneProvider(fakeGeneStats{})

	reporter.sendRegistration()
	reporter.sendHeartbeat()
	if err := reporter.ReportStatus("degraded", map[string]any{"reason": "low_battery"}); err != nil {
		t.Fatalf("ReportStatus: %v", err)
	}
	if err := reporter.ReportEvent("sensor_alert", map[string]any{"sensor": "do"}); err != nil {
		t.Fatalf("ReportEvent: %v", err)
	}

	wantPaths := []string{"/fleet/register", "/fleet/heartbeat", "/fleet/status", "/fleet/events"}
	if len(requests) != len(wantPaths) {
		t.Fatalf("requests = %d, want %d", len(requests), len(wantPaths))
	}
	for i, want := range wantPaths {
		if requests[i].Path != want {
			t.Fatalf("request %d path = %s, want %s", i, requests[i].Path, want)
		}
		if requests[i].Authorization != "Bearer fleet-token" {
			t.Fatalf("request %d auth = %q, want bearer token", i, requests[i].Authorization)
		}
		if requests[i].Payload["node_id"] != "edge-01" {
			t.Fatalf("request %d node_id = %v, want edge-01", i, requests[i].Payload["node_id"])
		}
		if _, ok := requests[i].Payload["timestamp"]; i > 0 && !ok {
			t.Fatalf("request %d missing timestamp", i)
		}
	}

	if requests[0].Payload["node_name"] != "Pond Edge 01" {
		t.Fatalf("registration node_name = %v, want Pond Edge 01", requests[0].Payload["node_name"])
	}
	if _, ok := requests[1].Payload["gene_stats"]; !ok {
		t.Fatal("heartbeat missing gene_stats")
	}
	if requests[2].Payload["status"] != "degraded" {
		t.Fatalf("status = %v, want degraded", requests[2].Payload["status"])
	}
	if requests[3].Payload["type"] != "sensor_alert" {
		t.Fatalf("event type = %v, want sensor_alert", requests[3].Payload["type"])
	}
}

func TestReporterReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadGateway)
	}))
	defer server.Close()

	reporter := NewReporter(Config{
		NodeID: "edge-01",
		Cloud: CloudConfig{
			Endpoint:          server.URL,
			HeartbeatInterval: time.Second,
		},
	})

	if err := reporter.ReportEvent("sensor_alert", nil); err == nil {
		t.Fatal("ReportEvent error = nil, want HTTP status error")
	}
}

func TestReporterUsesDefaultHeartbeatIntervalAndIdempotentStop(t *testing.T) {
	reporter := NewReporter(Config{})

	if reporter.config.Cloud.HeartbeatInterval != DefaultConfig().Cloud.HeartbeatInterval {
		t.Fatalf("interval = %s, want default %s", reporter.config.Cloud.HeartbeatInterval, DefaultConfig().Cloud.HeartbeatInterval)
	}

	reporter.Stop()
	reporter.Stop()
}
