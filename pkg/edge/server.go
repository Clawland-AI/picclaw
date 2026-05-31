package edge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
func NewServer(cfg Config, buses ...*bus.MessageBus) *Server {
	var msgBus *bus.MessageBus
	if len(buses) > 0 {
		msgBus = buses[0]
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
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /api/health", s.handleHealthz)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("GET /api/v1/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/command", s.handleCommand)
	s.mux.HandleFunc("POST /api/v1/command", s.handleCommand)
	s.mux.HandleFunc("POST /api/message", s.handleMessage)
	s.mux.HandleFunc("POST /api/v1/message", s.handleMessage)
}

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":        "ok",
		"agent":         "picclaw",
		"node_id":       s.config.NodeID,
		"node_name":     s.config.NodeName,
		"uptime":        time.Since(s.start).String(),
		"bus_connected": s.bus != nil,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"node_id":       s.config.NodeID,
		"node_name":     s.config.NodeName,
		"status":        "running",
		"uptime":        time.Since(s.start).String(),
		"cloud":         s.config.Cloud.Endpoint,
		"bus_connected": s.bus != nil,
		"endpoints": []string{
			"POST /api/command",
			"POST /api/message",
			"GET /api/status",
			"GET /api/health",
			"GET /healthz",
		},
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var cmd edgeCommand
	if err := decodeJSON(r, &cmd); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if cmd.Type == "" {
		writeError(w, http.StatusBadRequest, "command type is required")
		return
	}

	log.Printf("[edge] received command: type=%s id=%s", cmd.Type, cmd.CommandID)
	if s.bus != nil {
		s.bus.PublishInbound(cmd.toInbound(s.config.NodeID))
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":     "accepted",
		"command_id": cmd.CommandID,
		"command":    cmd.Type,
		"queued":     s.bus != nil,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var msg edgeMessage
	if err := decodeJSON(r, &msg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if msg.Content == "" {
		writeError(w, http.StatusBadRequest, "message content is required")
		return
	}

	log.Printf("[edge] received message: id=%s sender=%s", msg.MessageID, msg.SenderID)
	if s.bus != nil {
		s.bus.PublishInbound(msg.toInbound(s.config.NodeID))
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":     "accepted",
		"message_id": msg.MessageID,
		"queued":     s.bus != nil,
	})
}

type edgeCommand struct {
	CommandID  string            `json:"command_id"`
	Type       string            `json:"type"`
	Payload    map[string]any    `json:"payload"`
	SenderID   string            `json:"sender_id"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata"`
}

type edgeMessage struct {
	MessageID  string            `json:"message_id"`
	SenderID   string            `json:"sender_id"`
	ChatID     string            `json:"chat_id"`
	Content    string            `json:"content"`
	Media      []string          `json:"media"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata"`
}

func (cmd edgeCommand) toInbound(nodeID string) bus.InboundMessage {
	metadata := cloneMetadata(cmd.Metadata)
	metadata["edge_node_id"] = nodeID
	metadata["kind"] = "command"
	metadata["command_id"] = cmd.CommandID
	metadata["payload"] = jsonString(cmd.Payload)

	senderID := cmd.SenderID
	if senderID == "" {
		senderID = "edge-cloud"
	}

	return bus.InboundMessage{
		Channel:    "edge-api",
		SenderID:   senderID,
		ChatID:     firstNonEmpty(cmd.CommandID, nodeID),
		Content:    cmd.Type,
		SessionKey: cmd.SessionKey,
		Metadata:   metadata,
	}
}

func (msg edgeMessage) toInbound(nodeID string) bus.InboundMessage {
	metadata := cloneMetadata(msg.Metadata)
	metadata["edge_node_id"] = nodeID
	metadata["kind"] = "message"
	metadata["message_id"] = msg.MessageID

	senderID := msg.SenderID
	if senderID == "" {
		senderID = "edge-cloud"
	}

	return bus.InboundMessage{
		Channel:    "edge-api",
		SenderID:   senderID,
		ChatID:     firstNonEmpty(msg.ChatID, msg.MessageID, nodeID),
		Content:    msg.Content,
		Media:      msg.Media,
		SessionKey: msg.SessionKey,
		Metadata:   metadata,
	}
}

func decodeJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
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
	out := make(map[string]string, len(metadata)+4)
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
