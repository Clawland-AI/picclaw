package edge

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

func newTestServer() (*Server, *bus.MessageBus) {
	msgBus := bus.NewMessageBus()
	cfg := Config{
		Port:     9090,
		NodeID:   "edge-01",
		NodeName: "north-gate",
		Cloud: CloudConfig{
			Endpoint: "https://fleet.example.test",
		},
	}
	return NewServer(cfg, msgBus), msgBus
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return body
}

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

func TestHealthEndpointAliases(t *testing.T) {
	server, _ := newTestServer()

	for _, path := range []string{"/healthz", "/api/health"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, rec.Code)
		}
		body := decodeResponse(t, rec)
		if body["status"] != "ok" || body["node_id"] != "edge-01" || body["bus_connected"] != true {
			t.Fatalf("unexpected health body for %s: %#v", path, body)
		}
	}
}

func TestStatusEndpoint(t *testing.T) {
	server, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status returned %d", rec.Code)
	}
	body := decodeResponse(t, rec)
	if body["status"] != "running" || body["cloud"] != "https://fleet.example.test" {
		t.Fatalf("unexpected status body: %#v", body)
	}
}

func TestCommandEndpointPublishesToMessageBus(t *testing.T) {
	server, msgBus := newTestServer()
	payload := []byte(`{"command_id":"cmd-1","type":"restart_sensor","payload":{"sensor":"temp"},"sender_id":"fleet","session_key":"session-1","metadata":{"priority":"high"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("command returned %d: %s", rec.Code, rec.Body.String())
	}
	body := decodeResponse(t, rec)
	if body["queued"] != true || body["command"] != "restart_sensor" {
		t.Fatalf("unexpected command response: %#v", body)
	}

	msg := consumeInbound(t, msgBus)
	if msg.Channel != "edge-api" || msg.SenderID != "fleet" || msg.Content != "restart_sensor" {
		t.Fatalf("unexpected command bus message: %#v", msg)
	}
	if msg.Metadata["kind"] != "command" || msg.Metadata["command_id"] != "cmd-1" || msg.Metadata["priority"] != "high" {
		t.Fatalf("unexpected command metadata: %#v", msg.Metadata)
	}
}

func TestMessageEndpointPublishesToMessageBus(t *testing.T) {
	server, msgBus := newTestServer()
	payload := []byte(`{"message_id":"msg-1","channel":"operator","sender_id":"fleet","chat_id":"ops","content":"check pump","media":["image.jpg"],"metadata":{"priority":"high"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("message returned %d: %s", rec.Code, rec.Body.String())
	}
	body := decodeResponse(t, rec)
	if body["queued"] != true || body["message_id"] != "msg-1" {
		t.Fatalf("unexpected message response: %#v", body)
	}

	msg := consumeInbound(t, msgBus)
	if msg.Channel != "operator" || msg.SenderID != "fleet" || msg.ChatID != "ops" || msg.Content != "check pump" {
		t.Fatalf("unexpected message bus message: %#v", msg)
	}
	if len(msg.Media) != 1 || msg.Media[0] != "image.jpg" {
		t.Fatalf("unexpected media: %#v", msg.Media)
	}
	if msg.Metadata["kind"] != "message" || msg.Metadata["message_id"] != "msg-1" || msg.Metadata["priority"] != "high" {
		t.Fatalf("unexpected message metadata: %#v", msg.Metadata)
	}
}

func TestInvalidJSONReturnsBadRequest(t *testing.T) {
	server, _ := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewReader([]byte("{")))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid JSON returned %d", rec.Code)
	}
}

func TestCommandEndpointRejectsMissingType(t *testing.T) {
	server, _ := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewReader([]byte(`{"payload":{}}`)))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("missing command type returned %d", rec.Code)
	}
}
