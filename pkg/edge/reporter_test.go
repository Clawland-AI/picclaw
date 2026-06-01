package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type reporterRequest struct {
	Path          string
	Authorization string
	NodeIDHeader  string
	Payload       map[string]any
}

type testGeneProvider struct{}

func (testGeneProvider) GetStats() map[string]interface{} {
	return map[string]interface{}{"genes": 3, "capsules": 2}
}

func (testGeneProvider) GetHighConfidenceGenes(float64, int) []map[string]interface{} {
	return []map[string]interface{}{{"id": "gene-a"}}
}

func newReporterForTest(endpoint string) *Reporter {
	reporter := NewReporter(Config{
		NodeID:   "node-1",
		NodeName: "pond-edge",
		Cloud: CloudConfig{
			Endpoint:          endpoint,
			Token:             "secret-token",
			HeartbeatInterval: time.Second,
			StatusInterval:    2 * time.Second,
		},
	})
	reporter.started = time.Date(2026, 5, 31, 11, 59, 0, 0, time.UTC)
	reporter.now = func() time.Time {
		return time.Date(2026, 5, 31, 12, 0, 0, 0, time.UTC)
	}
	return reporter
}

func captureReporterServer(t *testing.T) (*httptest.Server, *[]reporterRequest) {
	t.Helper()
	records := []reporterRequest{}
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
		records = append(records, reporterRequest{
			Path:          r.URL.Path,
			Authorization: r.Header.Get("Authorization"),
			NodeIDHeader:  r.Header.Get("X-Clawland-Node-ID"),
			Payload:       payload,
		})
		w.WriteHeader(http.StatusAccepted)
	}))
	return server, &records
}

func TestRegisterPostsNodeMetadataWithHeaders(t *testing.T) {
	server, records := captureReporterServer(t)
	defer server.Close()
	reporter := newReporterForTest(server.URL + "/")

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
	if record.NodeIDHeader != "node-1" {
		t.Fatalf("node header = %q", record.NodeIDHeader)
	}
	if record.Payload["node_id"] != "node-1" || record.Payload["node_name"] != "pond-edge" {
		t.Fatalf("payload missing node metadata: %#v", record.Payload)
	}
	if record.Payload["timestamp"] != "2026-05-31T12:00:00Z" {
		t.Fatalf("timestamp = %#v", record.Payload["timestamp"])
	}
}

func TestSendHeartbeatIncludesGeneStatsAndUptime(t *testing.T) {
	server, records := captureReporterServer(t)
	defer server.Close()
	reporter := newReporterForTest(server.URL)
	reporter.SetGeneProvider(testGeneProvider{})

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
	if record.Payload["uptime_seconds"].(float64) != 60 {
		t.Fatalf("uptime = %#v", record.Payload["uptime_seconds"])
	}
	stats, ok := record.Payload["gene_stats"].(map[string]any)
	if !ok || stats["genes"].(float64) != 3 {
		t.Fatalf("gene_stats = %#v", record.Payload["gene_stats"])
	}
}

func TestReportStatusIncludesHighConfidenceGenes(t *testing.T) {
	server, records := captureReporterServer(t)
	defer server.Close()
	reporter := newReporterForTest(server.URL)
	reporter.SetGeneProvider(testGeneProvider{})

	if err := reporter.ReportStatus("degraded", map[string]any{"queue_depth": 4}); err != nil {
		t.Fatalf("ReportStatus() error = %v", err)
	}

	record := (*records)[0]
	if record.Path != "/fleet/status" {
		t.Fatalf("path = %s, want /fleet/status", record.Path)
	}
	if record.Payload["status"] != "degraded" {
		t.Fatalf("status payload = %#v", record.Payload)
	}
	if _, ok := record.Payload["high_confidence_genes"].([]any); !ok {
		t.Fatalf("high_confidence_genes missing: %#v", record.Payload)
	}
	data := record.Payload["data"].(map[string]any)
	if data["queue_depth"].(float64) != 4 {
		t.Fatalf("data = %#v", data)
	}
}

func TestReportEventAndPublishGene(t *testing.T) {
	server, records := captureReporterServer(t)
	defer server.Close()
	reporter := newReporterForTest(server.URL)

	if err := reporter.ReportEvent("sensor.alert", map[string]any{"sensor": "temp"}); err != nil {
		t.Fatalf("ReportEvent() error = %v", err)
	}
	if err := reporter.PublishGene(map[string]interface{}{"id": "gene-a"}); err != nil {
		t.Fatalf("PublishGene() error = %v", err)
	}

	if (*records)[0].Path != "/fleet/events" || (*records)[0].Payload["type"] != "sensor.alert" {
		t.Fatalf("event record = %#v", (*records)[0])
	}
	if (*records)[1].Path != "/fleet/genes/publish" {
		t.Fatalf("gene path = %s", (*records)[1].Path)
	}
}

func TestQueueEventRejectsEmptyTypeAndReportsFullQueue(t *testing.T) {
	reporter := newReporterForTest("http://example.test")
	if reporter.QueueEvent(" ", nil) {
		t.Fatal("QueueEvent(empty) = true, want false")
	}
	for i := 0; i < defaultEventBufferSize; i++ {
		if !reporter.QueueEvent("sensor.alert", nil) {
			t.Fatalf("QueueEvent(%d) = false, want true before full", i)
		}
	}
	if reporter.QueueEvent("sensor.alert", nil) {
		t.Fatal("QueueEvent(full) = true, want false")
	}
}

func TestPostReturnsFailureDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad token", http.StatusUnauthorized)
	}))
	defer server.Close()
	reporter := newReporterForTest(server.URL)

	err := reporter.ReportEvent("sensor.alert", nil)
	if err == nil {
		t.Fatal("ReportEvent() error = nil, want failure")
	}
	if !strings.Contains(err.Error(), "status 401") || !strings.Contains(err.Error(), "bad token") {
		t.Fatalf("error = %q, want status and response body", err.Error())
	}
}

func TestPostRejectsEmptyEndpoint(t *testing.T) {
	reporter := newReporterForTest("")
	if err := reporter.ReportEvent("sensor.alert", nil); err == nil || !strings.Contains(err.Error(), "cloud endpoint is empty") {
		t.Fatalf("ReportEvent() error = %v, want empty endpoint error", err)
	}
}

func TestReportEventRejectsEmptyType(t *testing.T) {
	reporter := newReporterForTest("http://example.test")
	if err := reporter.ReportEvent(" ", nil); err == nil {
		t.Fatal("ReportEvent() error = nil, want validation error")
	}
}

func TestStopIsIdempotent(t *testing.T) {
	reporter := newReporterForTest("http://example.test")
	reporter.Stop()
	reporter.Stop()
}

func TestNewReporterNormalizesIntervals(t *testing.T) {
	reporter := NewReporter(Config{})
	if reporter.config.Cloud.HeartbeatInterval != defaultHeartbeatInterval {
		t.Fatalf("heartbeat interval = %s", reporter.config.Cloud.HeartbeatInterval)
	}
	if reporter.config.Cloud.StatusInterval != defaultStatusInterval {
		t.Fatalf("status interval = %s", reporter.config.Cloud.StatusInterval)
	}
}
