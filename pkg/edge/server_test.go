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

func testServer() (*Server, *bus.MessageBus) {
	msgBus := bus.NewMessageBus()
	cfg := Config{
		Port:     9090,
		NodeID:   "edge-01",
		NodeName: "Edge 01",
		Cloud: CloudConfig{
			Endpoint: "http://fleet.local",
		},
	}
	return NewServer(cfg, msgBus), msgBus
}

func decodeJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode response JSON: %v", err)
	}
	return body
}

func TestHealthEndpoint(t *testing.T) {
	s, _ := testServer()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	s.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := decodeJSON(t, rr)
	if body["status"] != "ok" {
		t.Fatalf("expected ok status, got %#v", body["status"])
	}
	if body["node_id"] != "edge-01" {
		t.Fatalf("expected node id in health response, got %#v", body["node_id"])
	}
}

func TestStatusEndpoint(t *testing.T) {
	s, _ := testServer()
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rr := httptest.NewRecorder()

	s.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := decodeJSON(t, rr)
	if body["node_name"] != "Edge 01" {
		t.Fatalf("expected node name in status response, got %#v", body["node_name"])
	}
	if body["message_bus"] != "connected" {
		t.Fatalf("expected message bus status, got %#v", body["message_bus"])
	}
}

func TestCommandEndpointQueuesMessageBusCommand(t *testing.T) {
	s, msgBus := testServer()
	payload := []byte(`{"id":"cmd-1","type":"restart","payload":{"delay_seconds":5}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	s.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rr.Code, rr.Body.String())
	}
	body := decodeJSON(t, rr)
	if body["queued"] != true {
		t.Fatalf("expected queued response, got %#v", body["queued"])
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg, ok := msgBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected command on message bus")
	}
	if msg.Channel != "edge-api" {
		t.Fatalf("expected edge-api channel, got %s", msg.Channel)
	}
	if msg.Metadata["kind"] != "command" || msg.Metadata["command_type"] != "restart" {
		t.Fatalf("unexpected command metadata: %#v", msg.Metadata)
	}
	if msg.Content == "" {
		t.Fatal("expected command content")
	}
}

func TestMessageEndpointQueuesInboundMessage(t *testing.T) {
	s, msgBus := testServer()
	payload := []byte(`{"content":"check pond oxygen","sender_id":"operator","chat_id":"pond-01","metadata":{"priority":"high"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	s.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rr.Code, rr.Body.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg, ok := msgBus.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("expected inbound message on message bus")
	}
	if msg.Content != "check pond oxygen" {
		t.Fatalf("unexpected message content: %s", msg.Content)
	}
	if msg.Metadata["priority"] != "high" || msg.Metadata["source"] != "edge-http" {
		t.Fatalf("unexpected metadata: %#v", msg.Metadata)
	}
}

func TestCommandEndpointRejectsMissingType(t *testing.T) {
	s, _ := testServer()
	req := httptest.NewRequest(http.MethodPost, "/api/command", bytes.NewReader([]byte(`{"payload":{}}`)))
	rr := httptest.NewRecorder()

	s.mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
