package edge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

// Server is the Edge API HTTP server.
// It exposes endpoints for health checks, sensor data ingestion,
// and command reception from the upstream Fleet Manager.
type Server struct {
	config Config
	mux    *http.ServeMux
	start  time.Time
	bus    *bus.MessageBus
}

// NewServer creates a new Edge API Server.
func NewServer(cfg Config, messageBus ...*bus.MessageBus) *Server {
	var msgBus *bus.MessageBus
	if len(messageBus) > 0 {
		msgBus = messageBus[0]
	}

	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
		start:  time.Now(),
		bus:    msgBus,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	for _, route := range []string{"GET /healthz", "GET /api/health", "GET /api/v1/health"} {
		s.mux.HandleFunc(route, s.handleHealthz)
	}
	for _, route := range []string{"GET /api/status", "GET /api/v1/status"} {
		s.mux.HandleFunc(route, s.handleStatus)
	}
	for _, route := range []string{"POST /api/command", "POST /api/v1/command"} {
		s.mux.HandleFunc(route, s.handleCommand)
	}
	for _, route := range []string{"POST /api/message", "POST /api/v1/message"} {
		s.mux.HandleFunc(route, s.handleMessage)
	}
}

// ServeHTTP lets callers embed or test the Edge API server as a normal HTTP handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":        "ok",
		"agent":         "picclaw",
		"node_id":       s.config.NodeID,
		"node_name":     s.config.NodeName,
		"uptime":        time.Since(s.start).String(),
		"bus_connected": s.bus != nil,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"node_id":       s.config.NodeID,
		"node_name":     s.config.NodeName,
		"status":        "running",
		"uptime":        time.Since(s.start).String(),
		"cloud":         s.config.Cloud.Endpoint,
		"bus_connected": s.bus != nil,
		"endpoints": []string{
			"GET /healthz",
			"GET /api/health",
			"GET /api/status",
			"POST /api/command",
			"POST /api/message",
		},
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var req edgeCommandRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	req.CommandID = firstNonEmpty(req.CommandID, req.ID, fmt.Sprintf("cmd-%d", time.Now().UnixNano()))
	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		writeError(w, http.StatusBadRequest, "command type is required")
		return
	}

	if s.bus != nil {
		s.bus.PublishInbound(req.toInbound(s.config.NodeID))
	}

	log.Printf("[edge] queued command: id=%s type=%s", req.CommandID, req.Type)
	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":     "accepted",
		"command_id": req.CommandID,
		"command":    req.Type,
		"queued":     s.bus != nil,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var req edgeMessageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	req.MessageID = firstNonEmpty(req.MessageID, req.ID, fmt.Sprintf("msg-%d", time.Now().UnixNano()))
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "message content is required")
		return
	}

	if s.bus != nil {
		s.bus.PublishInbound(req.toInbound(s.config.NodeID))
	}

	log.Printf("[edge] queued message: id=%s sender=%s", req.MessageID, firstNonEmpty(req.SenderID, "fleet"))
	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":     "accepted",
		"message_id": req.MessageID,
		"queued":     s.bus != nil,
	})
}

type edgeCommandRequest struct {
	ID         string            `json:"id"`
	CommandID  string            `json:"command_id"`
	Type       string            `json:"type"`
	Payload    map[string]any    `json:"payload"`
	SenderID   string            `json:"sender_id"`
	ChatID     string            `json:"chat_id"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata"`
}

type edgeMessageRequest struct {
	ID         string            `json:"id"`
	MessageID  string            `json:"message_id"`
	Channel    string            `json:"channel"`
	SenderID   string            `json:"sender_id"`
	ChatID     string            `json:"chat_id"`
	Content    string            `json:"content"`
	Media      []string          `json:"media,omitempty"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

func (req edgeCommandRequest) toInbound(nodeID string) bus.InboundMessage {
	metadata := cloneMetadata(req.Metadata)
	metadata["kind"] = "command"
	metadata["source"] = "edge-http"
	metadata["edge_node_id"] = nodeID
	metadata["command_id"] = req.CommandID
	metadata["command_type"] = req.Type
	metadata["payload"] = jsonString(req.Payload)

	content := jsonString(map[string]any{
		"id":      req.CommandID,
		"type":    req.Type,
		"payload": req.Payload,
	})

	return bus.InboundMessage{
		Channel:    "edge-api",
		SenderID:   firstNonEmpty(req.SenderID, "fleet"),
		ChatID:     firstNonEmpty(req.ChatID, req.CommandID, nodeID),
		Content:    content,
		SessionKey: firstNonEmpty(req.SessionKey, "edge-command:"+req.CommandID),
		Metadata:   metadata,
	}
}

func (req edgeMessageRequest) toInbound(nodeID string) bus.InboundMessage {
	metadata := cloneMetadata(req.Metadata)
	metadata["kind"] = "message"
	metadata["source"] = "edge-http"
	metadata["edge_node_id"] = nodeID
	metadata["message_id"] = req.MessageID

	return bus.InboundMessage{
		Channel:    firstNonEmpty(req.Channel, "edge-api"),
		SenderID:   firstNonEmpty(req.SenderID, "fleet"),
		ChatID:     firstNonEmpty(req.ChatID, req.MessageID, nodeID),
		Content:    req.Content,
		Media:      req.Media,
		SessionKey: firstNonEmpty(req.SessionKey, "edge-message:"+firstNonEmpty(req.ChatID, req.MessageID, nodeID)),
		Metadata:   metadata,
	}
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func cloneMetadata(metadata map[string]string) map[string]string {
	out := make(map[string]string, len(metadata)+5)
	for key, value := range metadata {
		out[key] = value
	}
	return out
}

func jsonString(value any) string {
	if value == nil {
		return "{}"
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
