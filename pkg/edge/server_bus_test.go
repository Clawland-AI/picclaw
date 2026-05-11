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

func TestCommandEndpointPublishesToMessageBus(t *testing.T) {
	msgBus := bus.NewMessageBus()
	server := NewServerWithBus(Config{NodeID: "rack-a1", NodeName: "Rack A1"}, msgBus)

	payload := []byte(`{"type":"reboot","payload":{"delay":"10s"},"request_id":"cmd-123"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewReader(payload))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, response.Code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, ok := msgBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected command to be published to message bus")
	}
	if msg.Channel != "edge-command" {
		t.Fatalf("expected edge-command channel, got %q", msg.Channel)
	}
	if msg.Content != "reboot" {
		t.Fatalf("expected command content reboot, got %q", msg.Content)
	}
	if msg.SessionKey != "edge:rack-a1:cmd-123" {
		t.Fatalf("expected command session key, got %q", msg.SessionKey)
	}
	if msg.Metadata["payload"] != `{"delay":"10s"}` {
		t.Fatalf("expected serialized payload metadata, got %#v", msg.Metadata["payload"])
	}
}

func TestMessageEndpointPublishesInboundMessage(t *testing.T) {
	msgBus := bus.NewMessageBus()
	server := NewServerWithBus(Config{NodeID: "rack-a1", NodeName: "Rack A1"}, msgBus)

	body := map[string]any{
		"channel":     "http",
		"sender_id":   "fleet-manager",
		"chat_id":     "rack-a1",
		"content":     "collect temperature",
		"session_key": "ops:night-shift",
		"metadata": map[string]string{
			"priority": "high",
		},
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(encoded))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, response.Code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, ok := msgBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected message to be published to message bus")
	}
	if msg.Channel != "http" || msg.SenderID != "fleet-manager" || msg.ChatID != "rack-a1" {
		t.Fatalf("unexpected message routing: %#v", msg)
	}
	if msg.Content != "collect temperature" {
		t.Fatalf("expected content, got %q", msg.Content)
	}
	if msg.SessionKey != "ops:night-shift" {
		t.Fatalf("expected session key, got %q", msg.SessionKey)
	}
	if msg.Metadata["priority"] != "high" {
		t.Fatalf("expected metadata to round trip, got %#v", msg.Metadata)
	}
}

func TestMessageEndpointRequiresContent(t *testing.T) {
	server := NewServerWithBus(Config{NodeID: "rack-a1"}, bus.NewMessageBus())

	request := httptest.NewRequest(http.MethodPost, "/api/message", bytes.NewReader([]byte(`{"channel":"http"}`)))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}
