package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestHealthzReturnsNodeRuntimeMetadata(t *testing.T) {
	server := NewServer(Config{
		Port:      9090,
		NodeID:    "rack-a1",
		NodeName:  "Datacenter Rack A1",
		Version:   "test-version",
		BuildTime: "2026-05-11T00:00:00Z",
		AgentName: "picclaw-test",
	})

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected application/json content type, got %q", contentType)
	}

	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	expected := map[string]string{
		"status":     "ok",
		"version":    "test-version",
		"build_time": "2026-05-11T00:00:00Z",
		"agent_name": "picclaw-test",
		"node_id":    "rack-a1",
		"node_name":  "Datacenter Rack A1",
		"go_version": runtime.Version(),
	}
	for key, value := range expected {
		if body[key] != value {
			t.Fatalf("expected %s %q, got %#v", key, value, body[key])
		}
	}
	if _, ok := body["uptime"].(string); !ok {
		t.Fatalf("expected uptime string, got %#v", body["uptime"])
	}
}
