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

func newTestServer(messageBus *bus.MessageBus) *Server {
	return NewServer(Config{NodeID: "edge-01", NodeName: "PicClaw Lab", Cloud: CloudConfig{Endpoint: "https://fleet.example.test"}}, messageBus)
}

func TestHealthAndStatusRoutes(t *testing.T) {
	srv := newTestServer(bus.NewMessageBus())
	for _, path := range []string{"/healthz", "/api/health", "/api/v1/health"} {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, rec.Code)
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body["status"] != "ok" || body["node_id"] != "edge-01" || body["bus_connected"] != true {
			t.Fatalf("bad health body for %s: %#v", path, body)
		}
	}

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/status", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status returned %d", rec.Code)
	}
	var status map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
		t.Fatal(err)
	}
	if status["cloud"] != "https://fleet.example.test" || status["status"] != "running" {
		t.Fatalf("bad status: %#v", status)
	}
}

func TestCommandRoutesPublishToMessageBus(t *testing.T) {
	for _, path := range []string{"/api/command", "/api/v1/command"} {
		msgBus := bus.NewMessageBus()
		srv := newTestServer(msgBus)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"id":"legacy-id","command_id":"cmd-1","type":"restart_sensor","payload":{"sensor":"temp"},"sender_id":"fleet","session_key":"session-1"}`)))
		if rec.Code != http.StatusAccepted {
			t.Fatalf("%s returned %d: %s", path, rec.Code, rec.Body.String())
		}
		msg := consumeInbound(t, msgBus)
		if msg.Channel != "edge-api" || msg.ChatID != "cmd-1" || msg.SessionKey != "session-1" {
			t.Fatalf("bad command message: %#v", msg)
		}
		if msg.Metadata["kind"] != "command" || msg.Metadata["command_type"] != "restart_sensor" || msg.Metadata["payload"] != `{"sensor":"temp"}` {
			t.Fatalf("bad command metadata: %#v", msg.Metadata)
		}
		if !strings.Contains(msg.Content, "restart_sensor") {
			t.Fatalf("content should include command details: %s", msg.Content)
		}
	}
}

func TestMessageRoutesPublishToMessageBus(t *testing.T) {
	for _, path := range []string{"/api/message", "/api/v1/message"} {
		msgBus := bus.NewMessageBus()
		srv := newTestServer(msgBus)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"message_id":"msg-1","channel":"ops","sender_id":"fleet","chat_id":"pond","content":"check oxygen","media":["image.jpg"],"metadata":{"priority":"high"}}`)))
		if rec.Code != http.StatusAccepted {
			t.Fatalf("%s returned %d: %s", path, rec.Code, rec.Body.String())
		}
		msg := consumeInbound(t, msgBus)
		if msg.Channel != "ops" || msg.Content != "check oxygen" || msg.ChatID != "pond" {
			t.Fatalf("bad message: %#v", msg)
		}
		if len(msg.Media) != 1 || msg.Media[0] != "image.jpg" {
			t.Fatalf("bad media: %#v", msg.Media)
		}
		if msg.Metadata["priority"] != "high" || msg.Metadata["kind"] != "message" || msg.Metadata["source"] != "edge-http" {
			t.Fatalf("bad message metadata: %#v", msg.Metadata)
		}
	}
}

func TestInvalidRequestsReturnBadRequest(t *testing.T) {
	srv := newTestServer(bus.NewMessageBus())
	cases := []struct{ method, path, body string }{
		{http.MethodPost, "/api/command", "{"},
		{http.MethodPost, "/api/command", `{"payload":{}}`},
		{http.MethodPost, "/api/message", "{"},
		{http.MethodPost, "/api/message", `{"content":""}`},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s %s returned %d", tc.method, tc.path, rec.Code)
		}
	}
}

func TestServerWithoutMessageBusStillAccepts(t *testing.T) {
	srv := newTestServer(nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/command", strings.NewReader(`{"type":"status"}`)))
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected accepted without bus, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["queued"] != false {
		t.Fatalf("expected queued=false without bus: %#v", body)
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
