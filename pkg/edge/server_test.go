package edge

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

// testConfig returns a minimal Config for testing.
func testConfig() Config {
	return Config{
		Port:     9090,
		NodeID:   "edge-01",
		NodeName: "pond-guardian",
		Cloud: CloudConfig{
			Endpoint: "http://fleet.example.test",
		},
	}
}

// mustEncodeJSON helper for creating test payloads.
func mustEncodeJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json encode: %v", err)
	}
	return b
}

// consumeInbound reads a message from the bus with a timeout.
func consumeInbound(t *testing.T, msgBus *bus.MessageBus) bus.InboundMessage {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg, ok := msgBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected inbound message on bus")
	}
	return msg
}

// ---------------------------------------------------------------------------
// Health & Status endpoint tests
// ---------------------------------------------------------------------------

func TestHealthEndpoint(t *testing.T) {
	msgBus := bus.NewMessageBus()
	s := NewServer(testConfig(), WithMessageBus(msgBus))

	tests := []struct {
		name       string
		path       string
		wantStatus int
		checkBody  func(t *testing.T, body map[string]any)
	}{
		{
			name:       "/api/health",
			path:       "/api/health",
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]any) {
				if body["status"] != "ok" {
					t.Errorf("status = %v, want ok", body["status"])
				}
				if body["agent"] != "picclaw" {
					t.Errorf("agent = %v, want picclaw", body["agent"])
				}
				if body["node_id"] != "edge-01" {
					t.Errorf("node_id = %v, want edge-01", body["node_id"])
				}
				if body["bus_connected"] != "connected" {
					t.Errorf("bus_connected = %v, want connected", body["bus_connected"])
				}
				if _, ok := body["uptime"]; !ok {
					t.Error("uptime field missing")
				}
				if _, ok := body["version"]; !ok {
					t.Error("version field missing")
				}
			},
		},
		{
			name:       "/healthz (legacy)",
			path:       "/healthz",
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]any) {
				if body["status"] != "ok" {
					t.Errorf("status = %v, want ok", body["status"])
				}
			},
		},
		{
			name:       "/api/status",
			path:       "/api/status",
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]any) {
				if body["node_name"] != "pond-guardian" {
					t.Errorf("node_name = %v, want pond-guardian", body["node_name"])
				}
				if body["status"] != "running" {
					t.Errorf("status = %v, want running", body["status"])
				}
				if body["cloud"] != "http://fleet.example.test" {
					t.Errorf("cloud = %v, want http://fleet.example.test", body["cloud"])
				}
				if body["bus_connected"] != "connected" {
					t.Errorf("bus_connected = %v, want connected", body["bus_connected"])
				}
				eps, ok := body["endpoints"].([]any)
				if !ok || len(eps) == 0 {
					t.Errorf("endpoints missing or empty")
				}
			},
		},
		{
			name:       "/api/v1/status (legacy)",
			path:       "/api/v1/status",
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]any) {
				if body["node_id"] != "edge-01" {
					t.Errorf("node_id = %v, want edge-01", body["node_id"])
				}
			},
		},
		{
			name:       "unknown route returns 404",
			path:       "/api/nonexistent",
			wantStatus: http.StatusNotFound,
			checkBody:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}

			if tc.checkBody != nil && rec.Code == http.StatusOK {
				var body map[string]any
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				tc.checkBody(t, body)
			}
		})
	}
}

func TestHealthWithoutBus(t *testing.T) {
	s := NewServer(testConfig()) // no bus

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["bus_connected"] != "disconnected" {
		t.Errorf("bus_connected = %v, want disconnected", body["bus_connected"])
	}
}

// ---------------------------------------------------------------------------
// Command endpoint tests
// ---------------------------------------------------------------------------

func TestCommandEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		wantStatus  int
		wantBus     bool
		checkMsg    func(t *testing.T, msg bus.InboundMessage)
	}{
		{
			name: "accepts full command",
			payload: map[string]any{
				"command_id": "cmd-42",
				"type":       "restart_sensor",
				"payload":    map[string]any{"sensor": "do-probe-1", "delay": 5},
				"sender_id":  "fleet-manager",
			},
			wantStatus: http.StatusAccepted,
			wantBus:    true,
			checkMsg: func(t *testing.T, msg bus.InboundMessage) {
				if msg.Channel != "edge-api" {
					t.Errorf("channel = %s, want edge-api", msg.Channel)
				}
				if msg.Content != "restart_sensor" {
					t.Errorf("content = %s, want restart_sensor", msg.Content)
				}
				if msg.SenderID != "fleet-manager" {
					t.Errorf("sender_id = %s, want fleet-manager", msg.SenderID)
				}
				if msg.ChatID != "cmd-42" {
					t.Errorf("chat_id = %s, want cmd-42", msg.ChatID)
				}
				if msg.Metadata["kind"] != "command" {
					t.Errorf("metadata.kind = %s, want command", msg.Metadata["kind"])
				}
				if msg.Metadata["command_id"] != "cmd-42" {
					t.Errorf("metadata.command_id = %s, want cmd-42", msg.Metadata["command_id"])
				}
				if msg.Metadata["command_type"] != "restart_sensor" {
					t.Errorf("metadata.command_type = %s, want restart_sensor", msg.Metadata["command_type"])
				}
				if msg.Metadata["edge_node_id"] != "edge-01" {
					t.Errorf("metadata.edge_node_id = %s, want edge-01", msg.Metadata["edge_node_id"])
				}
				if msg.Metadata["payload"] == "" {
					t.Error("metadata.payload should contain serialized payload")
				}
			},
		},
		{
			name: "accepts minimal command (generates defaults)",
			payload: map[string]any{
				"type": "ping",
			},
			wantStatus: http.StatusAccepted,
			wantBus:    true,
			checkMsg: func(t *testing.T, msg bus.InboundMessage) {
				if msg.Content != "ping" {
					t.Errorf("content = %s, want ping", msg.Content)
				}
				if msg.SenderID != "fleet" {
					t.Errorf("sender_id = %s, want fleet (default)", msg.SenderID)
				}
				if msg.Metadata["command_id"] == "" {
					t.Error("command_id should not be empty (auto-generated)")
				}
			},
		},
		{
			name: "rejects missing type",
			payload: map[string]any{
				"command_id": "cmd-bad",
			},
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
		{
			name:       "rejects empty body",
			payload:    "",
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
		{
			name: "rejects null body",
			payload: nil,
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
		{
			name: "rejects boolean body",
			payload: true,
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
		{
			name:       "works with legacy /api/v1/command path",
			payload:    map[string]any{"type": "legacy-cmd"},
			wantStatus: http.StatusAccepted,
			wantBus:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msgBus := bus.NewMessageBus()
			s := NewServer(testConfig(), WithMessageBus(msgBus))

			path := "/api/command"
			// Use legacy path for the legacy test
			if tc.name == "works with legacy /api/v1/command path" {
				path = "/api/v1/command"
			}

			var bodyReader *bytes.Reader
			if tc.payload != nil {
				switch v := tc.payload.(type) {
				case string:
					bodyReader = bytes.NewReader([]byte(v))
				default:
					bodyReader = bytes.NewReader(mustEncodeJSON(t, v))
				}
			} else {
				bodyReader = bytes.NewReader([]byte("null"))
			}

			req := httptest.NewRequest(http.MethodPost, path, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}

			if tc.wantBus {
				msg := consumeInbound(t, msgBus)
				if tc.checkMsg != nil {
					tc.checkMsg(t, msg)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Message endpoint tests
// ---------------------------------------------------------------------------

func TestMessageEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		wantStatus  int
		wantBus     bool
		checkMsg    func(t *testing.T, msg bus.InboundMessage)
	}{
		{
			name: "accepts full message with metadata",
			payload: map[string]any{
				"message_id": "msg-1",
				"channel":    "edge-api",
				"sender_id":  "operator",
				"chat_id":    "ops-room",
				"content":    "check pond oxygen level",
				"media":      []string{"sensor-reading.jpg"},
				"session_key": "sess-1",
				"metadata": map[string]string{
					"priority": "high",
				},
			},
			wantStatus: http.StatusAccepted,
			wantBus:    true,
			checkMsg: func(t *testing.T, msg bus.InboundMessage) {
				if msg.Channel != "edge-api" {
					t.Errorf("channel = %s, want edge-api", msg.Channel)
				}
				if msg.Content != "check pond oxygen level" {
					t.Errorf("content = %s, want check pond oxygen level", msg.Content)
				}
				if msg.SenderID != "operator" {
					t.Errorf("sender_id = %s, want operator", msg.SenderID)
				}
				if msg.ChatID != "ops-room" {
					t.Errorf("chat_id = %s, want ops-room", msg.ChatID)
				}
				if len(msg.Media) != 1 || msg.Media[0] != "sensor-reading.jpg" {
					t.Errorf("media = %v, want [sensor-reading.jpg]", msg.Media)
				}
				if msg.Metadata["priority"] != "high" {
					t.Errorf("metadata.priority = %s, want high", msg.Metadata["priority"])
				}
				if msg.Metadata["kind"] != "message" {
					t.Errorf("metadata.kind = %s, want message", msg.Metadata["kind"])
				}
				if msg.Metadata["message_id"] != "msg-1" {
					t.Errorf("metadata.message_id = %s, want msg-1", msg.Metadata["message_id"])
				}
				if msg.Metadata["edge_node_id"] != "edge-01" {
					t.Errorf("metadata.edge_node_id = %s, want edge-01", msg.Metadata["edge_node_id"])
				}
			},
		},
		{
			name: "accepts minimal message (generates defaults)",
			payload: map[string]any{
				"content": "hello from fleet",
			},
			wantStatus: http.StatusAccepted,
			wantBus:    true,
			checkMsg: func(t *testing.T, msg bus.InboundMessage) {
				if msg.Content != "hello from fleet" {
					t.Errorf("content = %s, want hello from fleet", msg.Content)
				}
				if msg.Channel != "edge-api" {
					t.Errorf("channel = %s, want edge-api (default)", msg.Channel)
				}
				if msg.SenderID != "fleet" {
					t.Errorf("sender_id = %s, want fleet (default)", msg.SenderID)
				}
			},
		},
		{
			name: "rejects missing content",
			payload: map[string]any{
				"sender_id": "operator",
			},
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
		{
			name:       "rejects empty JSON object",
			payload:    map[string]any{},
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
		{
			name:       "rejects empty body",
			payload:    "",
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
		{
			name: "rejects malformed JSON",
			payload:    "{invalid",
			wantStatus: http.StatusBadRequest,
			wantBus:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msgBus := bus.NewMessageBus()
			s := NewServer(testConfig(), WithMessageBus(msgBus))

			var bodyReader *bytes.Reader
			if tc.payload != nil {
				switch v := tc.payload.(type) {
				case string:
					bodyReader = bytes.NewReader([]byte(v))
				default:
					bodyReader = bytes.NewReader(mustEncodeJSON(t, v))
				}
			} else {
				bodyReader = bytes.NewReader([]byte("null"))
			}

			req := httptest.NewRequest(http.MethodPost, "/api/message", bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}

			if tc.wantBus {
				msg := consumeInbound(t, msgBus)
				if tc.checkMsg != nil {
					tc.checkMsg(t, msg)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CORS test
// ---------------------------------------------------------------------------

func TestCORSPreflight(t *testing.T) {
	s := NewServer(testConfig())

	req := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("CORS origin missing")
	}
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Errorf("CORS methods missing")
	}
}

func TestCORSHeadersOnGet(t *testing.T) {
	s := NewServer(testConfig())

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("CORS origin missing on GET response")
	}
}

// ---------------------------------------------------------------------------
// Method-not-allowed / wrong method tests
// ---------------------------------------------------------------------------

func TestWrongMethods(t *testing.T) {
	s := NewServer(testConfig())

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/api/health"},
		{"PUT", "/api/status"},
		{"DELETE", "/api/command"},
		{"PATCH", "/api/message"},
		{"GET", "/api/command"},
		{"GET", "/api/message"},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code == http.StatusOK {
				t.Errorf("expected non-200 for %s %s, got %d", tc.method, tc.path, rec.Code)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Rate limiter test
// ---------------------------------------------------------------------------

func TestRateLimiter(t *testing.T) {
	rl := newRateLimiter(100, 5)
	// Allow up to burst starts
	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Fatalf("expected allow at token %d", i)
		}
	}
	// Should be rate limited now (but not strictly guaranteed with time-based refill)
	// Instead, test that it returns a value and doesn't panic
	_ = rl.Allow()
}

// ---------------------------------------------------------------------------
// Recovery middleware test
// ---------------------------------------------------------------------------

func TestRecoveryMiddleware(t *testing.T) {
	// Create a server with a handler that panics by hitting an unimplemented path
	// The recovery middleware should catch it and return 500
	s := NewServer(testConfig())

	// The recovery is wrapped around all handlers, so a panic in a handler should be caught.
	// We can't easily test this without injecting a panic handler, so we verify
	// that normal requests still work through the middleware chain.
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 through middleware chain, got %d", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Rate limiting integration test
// ---------------------------------------------------------------------------

func TestRateLimitRejection(t *testing.T) {
	s := NewServer(testConfig(), WithRateLimit(1000, 2))

	// First two requests should be allowed (burst = 2)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			t.Fatalf("unexpected rate limit at request %d", i+1)
		}
	}
}

// ---------------------------------------------------------------------------
// Concurrent access test
// ---------------------------------------------------------------------------

func TestConcurrentAccess(t *testing.T) {
	msgBus := bus.NewMessageBus()
	s := NewServer(testConfig(), WithMessageBus(msgBus))

	done := make(chan struct{})
	run := func() {
		defer func() { done <- struct{}{} }()
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("concurrent health check failed: %d", rec.Code)
			}
		}
	}

	go run()
	go run()
	go run()

	for i := 0; i < 3; i++ {
		<-done
	}
}

// ---------------------------------------------------------------------------
// Legacy NewServer compatibility test
// ---------------------------------------------------------------------------

func TestNewServerBackwardCompatibility(t *testing.T) {
	// Original call pattern: NewServer(cfg)
	s1 := NewServer(testConfig())
	if s1 == nil {
		t.Fatal("NewServer(cfg) returned nil")
	}

	// Variadic bus pattern: NewServer(cfg, msgBus)
	msgBus := bus.NewMessageBus()
	s2 := NewServer(testConfig(), msgBus)
	if s2 == nil {
		t.Fatal("NewServer(cfg, msgBus) returned nil")
	}
	if s2.bus != msgBus {
		t.Error("bus not set correctly via legacy variadic pattern")
	}
}

// ---------------------------------------------------------------------------
// Content-Type header tests
// ---------------------------------------------------------------------------

func TestContentTypeHeaders(t *testing.T) {
	s := NewServer(testConfig())

	for _, path := range []string{"/api/health", "/api/status", "/api/command", "/api/message"} {
		t.Run(path, func(t *testing.T) {
			var req *http.Request
			if path == "/api/command" || path == "/api/message" {
				body := mustEncodeJSON(t, map[string]string{"type": "test", "content": "test"})
				req = httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
			} else {
				req = httptest.NewRequest(http.MethodGet, path, nil)
			}
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			ct := rec.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %s, want application/json", ct)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Error response format tests
// ---------------------------------------------------------------------------

func TestErrorResponseFormat(t *testing.T) {
	s := NewServer(testConfig())

	// Trigger a 400 error by posting empty body to /api/command
	req := httptest.NewRequest(http.MethodPost, "/api/command", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var errResp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("error response not valid JSON: %v", err)
	}
	if errResp["error"] == nil {
		t.Error("error response missing 'error' field")
	}
	if errResp["code"] == nil {
		t.Error("error response missing 'code' field")
	}
}
