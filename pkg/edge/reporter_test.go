package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type requestRecord struct {
	Path          string
	Authorization string
	Payload       map[string]any
}

type fakeGeneProvider struct{}

func (fakeGeneProvider) GetStats() map[string]interface{} {
	return map[string]interface{}{"genes": 3, "capsules": 2}
}

func (fakeGeneProvider) GetHighConfidenceGenes(float64, int) []map[string]interface{} {
	return []map[string]interface{}{{"id": "gene-a"}}
}

func newTestReporter(serverURL string) *Reporter {
	reporter := NewReporter(Config{
		NodeID:   "node-1",
		NodeName: "pond-edge",
		Cloud: CloudConfig{
			Endpoint:          serverURL,
			Token:             "secret-token",
			HeartbeatInterval: time.Second,
		},
	})
	reporter.now = func() time.Time {
		return time.Date(2026, 5, 31, 12, 0, 0, 0, time.UTC)
	}
	return reporter
}

func captureServer(t *testing.T) (*httptest.Server, *[]requestRecord) {
	t.Helper()
	records := []requestRecord{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("content type = %s, want application/json", got)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		records = append(records, requestRecord{
			Path:          r.URL.Path,
			Authorization: r.Header.Get("Authorization"),
			Payload:       payload,
		})
		w.WriteHeader(http.StatusAccepted)
	}))
	return server, &records
}

func TestRegisterPostsNodeMetadataWithAuthorization(t *testing.T) {
	server, records := captureServer(t)
	defer server.Close()

	reporter := newTestReporter(server.URL + "/")

	if err := reporter.Register(); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if len(*records) != 1 {
		t.Fatalf("records = %d, want 1", len(*records))
	}
	record := (*records)[0]
	if record.Path != "/fleet/register" {
		t.Fatalf("path = %s, want /fleet/register", record.Path)
	}
	if record.Authorization != "Bearer secret-token" {
		t.Fatalf("authorization = %q", record.Authorization)
	}
	if record.Payload["node_id"] != "node-1" || record.Payload["node_name"] != "pond-edge" {
		t.Fatalf("payload missing node metadata: %#v", record.Payload)
	}
	if record.Payload["timestamp"] != "2026-05-31T12:00:00Z" {
		t.Fatalf("timestamp = %#v", record.Payload["timestamp"])
	}
}

func TestSendHeartbeatIncludesGeneStats(t *testing.T) {
	server, records := captureServer(t)
	defer server.Close()

	reporter := newTestReporter(server.URL)
	reporter.SetGeneProvider(fakeGeneProvider{})

	if err := reporter.SendHeartbeat(); err != nil {
		t.Fatalf("SendHeartbeat() error = %v", err)
	}

	record := (*records)[0]
	if record.Path != "/fleet/heartbeat" {
		t.Fatalf("path = %s, want /fleet/heartbeat", record.Path)
	}
	if record.Payload["status"] != "online" {
		t.Fatalf("status = %#v", record.Payload["status"])
	}
	stats, ok := record.Payload["gene_stats"].(map[string]any)
	if !ok {
		t.Fatalf("gene_stats missing or wrong type: %#v", record.Payload["gene_stats"])
	}
	if stats["genes"].(float64) != 3 {
		t.Fatalf("gene_stats.genes = %#v", stats["genes"])
	}
}

func TestReportStatusAndEvent(t *testing.T) {
	server, records := captureServer(t)
	defer server.Close()

	reporter := newTestReporter(server.URL)

	if err := reporter.ReportStatus("degraded", map[string]any{"queue_depth": 4}); err != nil {
		t.Fatalf("ReportStatus() error = %v", err)
	}
	if err := reporter.ReportEvent("sensor.alert", map[string]any{"sensor": "temp"}); err != nil {
		t.Fatalf("ReportEvent() error = %v", err)
	}

	status := (*records)[0]
	if status.Path != "/fleet/status" {
		t.Fatalf("status path = %s", status.Path)
	}
	if status.Payload["status"] != "degraded" {
		t.Fatalf("status payload = %#v", status.Payload)
	}

	event := (*records)[1]
	if event.Path != "/fleet/events" {
		t.Fatalf("event path = %s", event.Path)
	}
	if event.Payload["type"] != "sensor.alert" {
		t.Fatalf("event type = %#v", event.Payload["type"])
	}
}

func TestPublishGenePostsGenePayload(t *testing.T) {
	server, records := captureServer(t)
	defer server.Close()

	reporter := newTestReporter(server.URL)

	if err := reporter.PublishGene(map[string]interface{}{"id": "gene-a"}); err != nil {
		t.Fatalf("PublishGene() error = %v", err)
	}

	record := (*records)[0]
	if record.Path != "/fleet/genes/publish" {
		t.Fatalf("path = %s, want /fleet/genes/publish", record.Path)
	}
	gene, ok := record.Payload["gene"].(map[string]any)
	if !ok || gene["id"] != "gene-a" {
		t.Fatalf("gene payload = %#v", record.Payload["gene"])
	}
}

func TestPostReturnsFailureDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad token", http.StatusUnauthorized)
	}))
	defer server.Close()

	reporter := newTestReporter(server.URL)

	err := reporter.ReportEvent("sensor.alert", nil)
	if err == nil {
		t.Fatal("ReportEvent() error = nil, want failure")
	}
	if !strings.Contains(err.Error(), "status 401") || !strings.Contains(err.Error(), "bad token") {
		t.Fatalf("error = %q, want status and response body", err.Error())
	}
}

func TestReportEventRejectsEmptyType(t *testing.T) {
	reporter := newTestReporter("http://example.test")

	if err := reporter.ReportEvent(" ", nil); err == nil {
		t.Fatal("ReportEvent() error = nil, want validation error")
	}
}

func TestStopIsIdempotent(t *testing.T) {
	reporter := newTestReporter("http://example.test")

	reporter.Stop()
	reporter.Stop()
}
