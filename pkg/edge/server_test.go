package edge

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

func testConfig() Config {
	return Config{
		Port:     9090,
		NodeID:   "edge-01",
		NodeName: "pond-guardian",
		Cloud:    CloudConfig{Endpoint: "http://fleet.example.test"},
	}
}

func consumeInbound(t *testing.T, mb *bus.MessageBus) bus.InboundMessage {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg, ok := mb.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected inbound message on bus")
	}
	return msg
}

func TestHealthAndStatus(t *testing.T) {
	s := NewServer(testConfig(), bus.NewMessageBus())
	for _, tc := range []struct {
		name       string
		path       string
		wantStatus int
		check      func(t *testing.T, body map[string]any)
	}{
		{"/api/health", "/api/health", 200, func(t *testing.T, b map[string]any) {
			if b["status"] != "ok" || b["agent"] != "picclaw" || b["node_id"] != "edge-01" || b["bus_connected"] != "connected" {
				t.Errorf("unexpected health body: %v", b)
			}
		}},
		{"/healthz legacy", "/healthz", 200, func(t *testing.T, b map[string]any) {
			if b["status"] != "ok" {
				t.Errorf("healthz status = %v", b["status"])
			}
		}},
		{"/api/status", "/api/status", 200, func(t *testing.T, b map[string]any) {
			if b["node_name"] != "pond-guardian" || b["status"] != "running" || b["bus_connected"] != "connected" {
				t.Errorf("unexpected status body: %v", b)
			}
		}},
		{"/api/v1/status legacy", "/api/v1/status", 200, func(t *testing.T, b map[string]any) {
			if b["node_id"] != "edge-01" {
				t.Errorf("node_id = %v", b["node_id"])
			}
		}},
		{"404", "/api/nonexistent", 404, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, httptest.NewRequest("GET", tc.path, nil))
			if rec.Code != tc.wantStatus {
				t.Fatalf("status %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.check != nil {
				var body map[string]any
				json.Unmarshal(rec.Body.Bytes(), &body)
				tc.check(t, body)
			}
		})
	}
}

func TestHealthWithoutBus(t *testing.T) {
	rec := httptest.NewRecorder()
	NewServer(testConfig(), nil).ServeHTTP(rec, httptest.NewRequest("GET", "/api/health", nil))
	var body map[string]any
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body["bus_connected"] != "disconnected" {
		t.Errorf("bus_connected = %v, want disconnected", body["bus_connected"])
	}
}

func TestCommandEndpoint(t *testing.T) {
	for _, tc := range []struct {
		name       string
		path       string
		body       string
		wantStatus int
		wantBus    bool
		check      func(t *testing.T, msg bus.InboundMessage)
	}{
		{"full command", "/api/command", `{"command_id":"c1","type":"reboot","payload":{"delay":5},"sender_id":"admin"}`, 202, true,
			func(t *testing.T, m bus.InboundMessage) {
				if m.Channel != "edge-api" || m.Content != "reboot" || m.SenderID != "admin" || m.ChatID != "c1" || m.Metadata["kind"] != "command" {
					t.Errorf("unexpected msg: %+v", m)
				}
			}},
		{"minimal", "/api/command", `{"type":"ping"}`, 202, true,
			func(t *testing.T, m bus.InboundMessage) {
				if m.Content != "ping" || m.SenderID != "fleet" || m.Metadata["command_id"] == "" {
					t.Errorf("unexpected msg: %+v", m)
				}
			}},
		{"missing type", "/api/command", `{"command_id":"x"}`, 400, false, nil},
		{"empty body", "/api/command", ``, 400, false, nil},
		{"legacy path", "/api/v1/command", `{"type":"legacy"}`, 202, true, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mb := bus.NewMessageBus()
			s := NewServer(testConfig(), mb)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", tc.path, bytes.NewReader([]byte(tc.body)))
			req.Header.Set("Content-Type", "application/json")
			s.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status %d, want %d: %s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantBus && tc.check != nil {
				tc.check(t, consumeInbound(t, mb))
			}
		})
	}
}

func TestMessageEndpoint(t *testing.T) {
	for _, tc := range []struct {
		name       string
		body       string
		wantStatus int
		wantBus    bool
		check      func(t *testing.T, msg bus.InboundMessage)
	}{
		{"full message", `{"message_id":"m1","channel":"edge-api","sender_id":"op","chat_id":"room","content":"hello","media":["img.jpg"],"metadata":{"prio":"high"}}`, 202, true,
			func(t *testing.T, m bus.InboundMessage) {
				if m.Channel != "edge-api" || m.Content != "hello" || m.SenderID != "op" || m.ChatID != "room" || len(m.Media) != 1 || m.Metadata["kind"] != "message" {
					t.Errorf("unexpected msg: %+v", m)
				}
			}},
		{"minimal", `{"content":"hi"}`, 202, true,
			func(t *testing.T, m bus.InboundMessage) {
				if m.Content != "hi" || m.Channel != "edge-api" || m.SenderID != "fleet" {
					t.Errorf("unexpected msg: %+v", m)
				}
			}},
		{"missing content", `{"sender_id":"op"}`, 400, false, nil},
		{"empty object", `{}`, 400, false, nil},
		{"empty body", ``, 400, false, nil},
		{"malformed JSON", `{invalid`, 400, false, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mb := bus.NewMessageBus()
			s := NewServer(testConfig(), mb)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/message", bytes.NewReader([]byte(tc.body)))
			req.Header.Set("Content-Type", "application/json")
			s.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status %d, want %d: %s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantBus && tc.check != nil {
				tc.check(t, consumeInbound(t, mb))
			}
		})
	}
}
