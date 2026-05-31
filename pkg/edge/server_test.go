package edge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

func testServer(msgBus *bus.MessageBus) *Server {
	return NewServer(Config{
		Port:     9090,
		NodeID:   "edge-01",
		NodeName: "PicClaw Lab",
		Cloud: CloudConfig{
			Endpoint: "https://fleet.example.test",
		},
	}, msgBus)
}

func TestHealthAndStatusEndpoints(t *testing.T) {
	msgBus := bus.NewMessageBus()
	server := testServer(msgBus)

	for _, path := range []string{"/healthz", "/api/health"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, rec.Code)
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode health response: %v", err)
		}
		if body["status"] != "ok" || body["node_id"] != "edge-01" || body["bus_connected"] != true {
			t.Fatalf("unexpected health response for %s: %#v", path, body)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status returned %d", rec.Code)
	}

	var status map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
		t.Fatalf("decode status response: %v", err)
	}
	if status["cloud"] != "https://fleet.example.test" || status["status"] != "running" {
		t.Fatalf("unexpected status response: %#v", status)
	}
}

func TestCommandEndpointPublishesToMessageBus(t *testing.T) {
	msgBus := bus.NewMessageBus()
	server := testServer(msgBus)

	body := `{"command_id":"cmd-1","type":"restart_sensor","payload":{"sensor":"temp"},"sender_id":"fleet","session_key":"session-1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/command", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("command returned %d: %s", rec.Code, rec.Body.String())
	}

	msg := consumeInbound(t, msgBus)
	if msg.Channel != "edge-api" || msg.Content != "restart_sensor" || msg.ChatID != "cmd-1" {
		t.Fatalf("unexpected inbound command: %#v", msg)
	}
	if msg.Metadata["kind"] != "command" || msg.Metadata["payload"] != `{"sensor":"temp"}` {
		t.Fatalf("unexpected command metadata: %#v", msg.Metadata)
	}
}

func TestMessageEndpointPublishesToMessageBus(t *testing.T) {
	msgBus := bus.NewMessageBus()
	server := testServer(msgBus)

	body := `{"message_id":"msg-1","sender_id":"fleet","chat_id":"ops","content":"check pump","media":["image.jpg"],"metadata":{"priority":"high"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/message", strings.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("message returned %d: %s", rec.Code, rec.Body.String())
	}

	msg := consumeInbound(t, msgBus)
	if msg.Channel != "edge-api" || msg.Content != "check pump" || msg.ChatID != "ops" {
		t.Fatalf("unexpected inbound message: %#v", msg)
	}
	if len(msg.Media) != 1 || msg.Media[0] != "image.jpg" {
		t.Fatalf("unexpected media: %#v", msg.Media)
	}
	if msg.Metadata["priority"] != "high" || msg.Metadata["kind"] != "message" {
		t.Fatalf("unexpected message metadata: %#v", msg.Metadata)
	}
}

func TestInvalidJSONReturnsBadRequest(t *testing.T) {
	server := testServer(bus.NewMessageBus())

	req := httptest.NewRequest(http.MethodPost, "/api/command", strings.NewReader("{"))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid JSON") {
		t.Fatalf("unexpected error response: %s", rec.Body.String())
	}
}

func consumeInbound(t *testing.T, msgBus *bus.MessageBus) bus.InboundMessage {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, ok := msgBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected inbound message")
	}
	return msg
}
