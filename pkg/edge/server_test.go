package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestHealthzIncludesRuntimeMetadata(t *testing.T) {
	server := NewServer(Config{
		NodeID:    "edge-01",
		NodeName:  "warehouse edge",
		AgentName: "picclaw-test",
		Version:   "test-version",
		BuildTime: "2026-05-31T12:00:00Z",
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	assertString(t, body, "status", "ok")
	assertString(t, body, "agent_name", "picclaw-test")
	assertString(t, body, "agent", "picclaw-test")
	assertString(t, body, "node_id", "edge-01")
	assertString(t, body, "version", "test-version")
	assertString(t, body, "build_time", "2026-05-31T12:00:00Z")
	assertString(t, body, "go_version", runtime.Version())

	if body["uptime"] == "" {
		t.Fatal("expected uptime to be present")
	}
}

func assertString(t *testing.T, body map[string]any, key string, want string) {
	t.Helper()
	if got, ok := body[key].(string); !ok || got != want {
		t.Fatalf("expected %s=%q, got %#v", key, want, body[key])
	}
}
