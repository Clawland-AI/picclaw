package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthzHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedOK     bool
	}{
		{
			name:           "GET request returns 200",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedOK:     true,
		},
		{
			name:           "POST request returns 405",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedOK:     false,
		},
		{
			name:           "PUT request returns 405",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/healthz", nil)
			w := httptest.NewRecorder()

			HealthzHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectedOK {
				var health HealthResponse
				if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if health.Status != "ok" {
					t.Errorf("expected status 'ok', got '%s'", health.Status)
				}

				if health.Version == "" {
					t.Error("expected non-empty version")
				}

				if health.AgentName == "" {
					t.Error("expected non-empty agent_name")
				}

				if health.GoVersion == "" {
					t.Error("expected non-empty go_version")
				}

				if health.Uptime == "" {
					t.Error("expected non-empty uptime")
				}
			}
		})
	}
}

func TestHealthzResponseFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	HealthzHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// Validate JSON structure
	expectedFields := map[string]bool{
		"status":     health.Status != "",
		"uptime":     health.Uptime != "",
		"version":    health.Version != "",
		"agent_name": health.AgentName != "",
		"go_version": health.GoVersion != "",
	}

	for field, exists := range expectedFields {
		if !exists {
			t.Errorf("expected field '%s' to be non-empty", field)
		}
	}
}

func TestNewHealthServer(t *testing.T) {
	addr := ":8080"
	srv := NewHealthServer(addr)

	if srv.Addr != addr {
		t.Errorf("expected address %s, got %s", addr, srv.Addr)
	}

	if srv.ReadTimeout != 10*time.Second {
		t.Errorf("expected ReadTimeout 10s, got %v", srv.ReadTimeout)
	}

	if srv.WriteTimeout != 10*time.Second {
		t.Errorf("expected WriteTimeout 10s, got %v", srv.WriteTimeout)
	}

	if srv.IdleTimeout != 30*time.Second {
		t.Errorf("expected IdleTimeout 30s, got %v", srv.IdleTimeout)
	}

	if srv.Handler == nil {
		t.Error("expected non-nil Handler")
	}
}

func TestHealthzUptime(t *testing.T) {
	// Reset StartTime to a known value
	oldStartTime := StartTime
	StartTime = time.Now().Add(-5 * time.Minute)
	defer func() { StartTime = oldStartTime }()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	HealthzHandler(w, req)

	var health HealthResponse
	json.NewDecoder(w.Result().Body).Decode(&health)

	// Parse the uptime duration
	duration, err := time.ParseDuration(health.Uptime)
	if err != nil {
		t.Fatalf("failed to parse uptime: %v", err)
	}

	// Should be approximately 5 minutes (allow 1 second tolerance)
	expected := 5 * time.Minute
	tolerance := 1 * time.Second

	if duration < expected-tolerance || duration > expected+tolerance {
		t.Errorf("expected uptime ~%v, got %v", expected, duration)
	}
}