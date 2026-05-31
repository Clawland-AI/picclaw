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
func NewServer(cfg Config, messageBus ...*bus.MessageBus) *Server {
	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
		start:  time.Now(),
	}
	if len(messageBus) > 0 {
		s.bus = messageBus[0]
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /api/health", s.handleHealthz)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/command", s.handleCommand)
	s.mux.HandleFunc("POST /api/message", s.handleMessage)

	// Backward-compatible v1 aliases kept for existing PicClaw configs.
	s.mux.HandleFunc("GET /api/v1/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/v1/command", s.handleCommand)
}

// ServeHTTP allows tests and embedded callers to use Server as an http.Handler.
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
		"status":    "ok",
		"agent":     "picclaw",
		"node_id":   s.config.NodeID,
		"node_name": s.config.NodeName,
		"uptime":    time.Since(s.start).String(),
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"node_id":   s.config.NodeID,
		"node_name": s.config.NodeName,
		"status":    "running",
		"uptime":    time.Since(s.start).String(),
		"cloud":     s.config.Cloud.Endpoint,
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var cmd commandRequest
	if err := decodeJSON(r, &cmd); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if cmd.Type == "" {
		http.Error(w, "missing command type", http.StatusBadRequest)
		return
	}

	payloadJSON, err := json.Marshal(cmd.Payload)
	if err != nil {
		http.Error(w, "invalid command payload", http.StatusBadRequest)
		return
	}

	inbound := bus.InboundMessage{
		Channel:    "edge",
		SenderID:   firstNonEmpty(cmd.SenderID, "fleet"),
		ChatID:     firstNonEmpty(cmd.CommandID, cmd.Type),
		Content:    cmd.Type,
		SessionKey: fmt.Sprintf("edge:%s", firstNonEmpty(cmd.CommandID, cmd.Type)),
		Metadata: map[string]string{
			"kind":         "command",
			"command_type": cmd.Type,
			"command_id":   cmd.CommandID,
			"payload_json": string(payloadJSON),
			"node_id":      s.config.NodeID,
		},
	}
	s.publishInbound(inbound)

	log.Printf("[edge] accepted command: type=%s id=%s", cmd.Type, cmd.CommandID)
	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":     "accepted",
		"command":    cmd.Type,
		"command_id": cmd.CommandID,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var msg edgeMessageRequest
	if err := decodeJSON(r, &msg); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if msg.Content == "" {
		http.Error(w, "missing message content", http.StatusBadRequest)
		return
	}

	inbound := bus.InboundMessage{
		Channel:    firstNonEmpty(msg.Channel, "edge"),
		SenderID:   firstNonEmpty(msg.SenderID, "fleet"),
		ChatID:     firstNonEmpty(msg.ChatID, "edge"),
		Content:    msg.Content,
		Media:      msg.Media,
		SessionKey: firstNonEmpty(msg.SessionKey, fmt.Sprintf("edge:%s", firstNonEmpty(msg.ChatID, "default"))),
		Metadata:   msg.Metadata,
	}
	if inbound.Metadata == nil {
		inbound.Metadata = map[string]string{}
	}
	inbound.Metadata["kind"] = "message"
	inbound.Metadata["node_id"] = s.config.NodeID

	s.publishInbound(inbound)
	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":      "accepted",
		"channel":     inbound.Channel,
		"session_key": inbound.SessionKey,
	})
}

func (s *Server) publishInbound(msg bus.InboundMessage) {
	if s.bus == nil {
		return
	}
	s.bus.PublishInbound(msg)
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

type commandRequest struct {
	CommandID string         `json:"command_id"`
	Type      string         `json:"type"`
	SenderID  string         `json:"sender_id"`
	Payload   map[string]any `json:"payload"`
}

type edgeMessageRequest struct {
	Channel    string            `json:"channel"`
	SenderID   string            `json:"sender_id"`
	ChatID     string            `json:"chat_id"`
	Content    string            `json:"content"`
	Media      []string          `json:"media,omitempty"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}
