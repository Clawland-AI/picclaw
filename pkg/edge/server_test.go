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

func testConfig() Config {
	cfg := DefaultConfig()
	cfg.NodeID = "edge-01"
	cfg.NodeName = "pond guardian"
	cfg.Cloud.Endpoint = "http://fleet.local"
	return cfg
}

func TestEdgeHealthAndStatusRoutes(t *testing.T) {
	server := NewServer(testConfig())

	for _, route := range []string{"/api/health", "/healthz", "/api/status", "/api/v1/status"} {
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()

		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s returned status %d", route, rec.Code)
		}

		var body map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("%s returned invalid JSON: %v", route, err)
		}
		if body["node_id"] != "edge-01" {
			t.Fatalf("%s node_id = %v, want edge-01", route, body["node_id"])
		}
	}
}

func TestEdgeCommandPublishesToMessageBus(t *testing.T) {
	messageBus := bus.NewMessageBus()
	server := NewServer(testConfig(), messageBus)

	req := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewBufferString(`{
		"command_id": "cmd-123",
		"type": "restart_sensor",
		"sender_id": "fleet-manager",
		"payload": {"sensor": "do-probe-1"}
	}`))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, ok := messageBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected command to be published to inbound message bus")
	}
	if msg.Channel != "edge" || msg.SenderID != "fleet-manager" || msg.Content != "restart_sensor" {
		t.Fatalf("unexpected message: %#v", msg)
	}
	if msg.Metadata["kind"] != "command" || msg.Metadata["command_id"] != "cmd-123" {
		t.Fatalf("unexpected command metadata: %#v", msg.Metadata)
	}
}

func TestEdgeMessagePublishesToMessageBus(t *testing.T) {
	messageBus := bus.NewMessageBus()
	server := NewServer(testConfig(), messageBus)

	req := httptest.NewRequest(http.MethodPost, "/api/message", bytes.NewBufferString(`{
		"channel": "edge-api",
		"sender_id": "operator",
		"chat_id": "ops-room",
		"content": "check pond oxygen",
		"metadata": {"priority": "high"}
	}`))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, ok := messageBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected message to be published to inbound message bus")
	}
	if msg.Channel != "edge-api" || msg.ChatID != "ops-room" || msg.Content != "check pond oxygen" {
		t.Fatalf("unexpected message: %#v", msg)
	}
	if msg.Metadata["kind"] != "message" || msg.Metadata["priority"] != "high" {
		t.Fatalf("unexpected message metadata: %#v", msg.Metadata)
	}
}

func TestEdgeRejectsInvalidJSON(t *testing.T) {
	server := NewServer(testConfig())

	req := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewBufferString(`{`))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
