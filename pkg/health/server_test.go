package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	cfg := Config{
		Port:      8080,
		Version:   "1.0.0",
		BuildTime: "2024-01-01",
		AgentName: "test-agent",
	}

	server := NewServer(cfg)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.port != 8080 {
		t.Errorf("Expected port 8080, got %d", server.port)
	}
	if server.version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", server.version)
	}
	if server.agentName != "test-agent" {
		t.Errorf("Expected agent name test-agent, got %s", server.agentName)
	}
}

func TestHealthzEndpoint(t *testing.T) {
	cfg := Config{
		Port:      8080,
		Version:   "1.0.0",
		BuildTime: "2024-01-01T00:00:00Z",
		AgentName: "picclaw",
	}

	server := NewServer(cfg)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	// Call handler
	server.handleHealthz(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check Content-Type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify required fields
	requiredFields := []string{"status", "uptime", "version", "agent_name", "go_version", "build_time"}
	for _, field := range requiredFields {
		if _, ok := response[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Verify status
	if response["status"] != "ok" {
		t.Errorf("Expected status ok, got %v", response["status"])
	}

	// Verify version
	if response["version"] != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %v", response["version"])
	}

	// Verify agent_name
	if response["agent_name"] != "picclaw" {
		t.Errorf("Expected agent_name picclaw, got %v", response["agent_name"])
	}

	// Verify go_version is not empty
	if goVersion, ok := response["go_version"].(string); !ok || goVersion == "" {
		t.Error("go_version is missing or empty")
	}
}

func TestHealthzMethodNotAllowed(t *testing.T) {
	cfg := Config{
		Port:      8080,
		Version:   "1.0.0",
		BuildTime: "2024-01-01",
		AgentName: "picclaw",
	}

	server := NewServer(cfg)

	// Test POST request (should be rejected)
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w := httptest.NewRecorder()

	server.handleHealthz(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHealthzUptime(t *testing.T) {
	cfg := Config{
		Port:      8080,
		Version:   "1.0.0",
		BuildTime: "2024-01-01",
		AgentName: "picclaw",
	}

	server := NewServer(cfg)

	// Wait a bit to ensure uptime > 0
	time.Sleep(100 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	server.handleHealthz(w, req)

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	uptime, ok := response["uptime"].(string)
	if !ok {
		t.Fatal("uptime field is not a string")
	}

	if uptime == "" || uptime == "0s" {
		t.Errorf("Expected uptime > 0, got %s", uptime)
	}
}

func TestHealthzMultipleRequests(t *testing.T) {
	cfg := Config{
		Port:      8080,
		Version:   "1.0.0",
		BuildTime: "2024-01-01",
		AgentName: "picclaw",
	}

	server := NewServer(cfg)

	// Send multiple requests
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()

		server.handleHealthz(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i, w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Request %d: Failed to decode response: %v", i, err)
		}

		if response["status"] != "ok" {
			t.Errorf("Request %d: Expected status ok, got %v", i, response["status"])
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestHealthzResponseFields(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		buildTime string
		agentName string
	}{
		{"Production", "1.0.0", "2024-01-01T00:00:00Z", "picclaw"},
		{"Development", "dev", "unknown", "picclaw-dev"},
		{"Empty version", "", "", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Port:      8080,
				Version:   tt.version,
				BuildTime: tt.buildTime,
				AgentName: tt.agentName,
			}

			server := NewServer(cfg)
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			w := httptest.NewRecorder()

			server.handleHealthz(w, req)

			var response map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response["version"] != tt.version {
				t.Errorf("Expected version %s, got %v", tt.version, response["version"])
			}

			if response["build_time"] != tt.buildTime {
				t.Errorf("Expected build_time %s, got %v", tt.buildTime, response["build_time"])
			}

			if response["agent_name"] != tt.agentName {
				t.Errorf("Expected agent_name %s, got %v", tt.agentName, response["agent_name"])
			}
		})
	}
}
