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
type Server struct {
	config Config
	mux    *http.ServeMux
	bus    *bus.MessageBus
	start  time.Time
}

// NewServer creates a new Edge API Server.
func NewServer(cfg Config, msgBus *bus.MessageBus) *Server {
	s := &Server{config: cfg, mux: http.NewServeMux(), start: time.Now(), bus: msgBus}
	s.mux.HandleFunc("GET /api/health", s.handleHealthz)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/command", s.handleCommand)
	s.mux.HandleFunc("POST /api/message", s.handleMessage)
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /api/v1/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/v1/command", s.handleCommand)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
	return http.ListenAndServe(addr, s)
}

// --- Handlers ---

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	busStatus := "disconnected"
	if s.bus != nil {
		busStatus = "connected"
	}
	writeJSON(w, 200, map[string]any{
		"status": "ok", "agent": "picclaw", "node_id": s.config.NodeID,
		"node_name": s.config.NodeName, "uptime": time.Since(s.start).String(),
		"version": "0.1.0", "bus_connected": busStatus,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	busInfo := "disconnected"
	if s.bus != nil {
		busInfo = "connected"
	}
	writeJSON(w, 200, map[string]any{
		"node_id": s.config.NodeID, "node_name": s.config.NodeName,
		"status": "running", "uptime": time.Since(s.start).String(),
		"cloud": s.config.Cloud.Endpoint, "bus_connected": busInfo,
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CommandID  string            `json:"command_id"`
		Type       string            `json:"type"`
		Payload    map[string]any    `json:"payload"`
		SenderID   string            `json:"sender_id"`
		SessionKey string            `json:"session_key"`
		Metadata   map[string]string `json:"metadata,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid JSON: " + err.Error(), "code": 400})
		return
	}
	if req.Type == "" {
		writeJSON(w, 400, map[string]any{"error": "command type is required", "code": 400})
		return
	}
	if req.CommandID == "" {
		req.CommandID = fmt.Sprintf("cmd-%d", time.Now().UnixNano())
	}
	if req.SenderID == "" {
		req.SenderID = "fleet"
	}
	if req.SessionKey == "" {
		req.SessionKey = "edge:cmd:" + req.CommandID
	}
	if req.Metadata == nil {
		req.Metadata = make(map[string]string)
	}
	req.Metadata["kind"] = "command"
	req.Metadata["command_id"] = req.CommandID
	req.Metadata["command_type"] = req.Type
	req.Metadata["edge_node_id"] = s.config.NodeID
	if req.Payload != nil {
		if b, err := json.Marshal(req.Payload); err == nil {
			req.Metadata["payload"] = string(b)
		}
	}
	if s.bus != nil {
		s.bus.PublishInbound(bus.InboundMessage{
			Channel: "edge-api", SenderID: req.SenderID, ChatID: req.CommandID,
			Content: req.Type, SessionKey: req.SessionKey, Metadata: req.Metadata,
		})
	}
	log.Printf("[edge] accepted command: type=%s id=%s node=%s", req.Type, req.CommandID, s.config.NodeID)
	writeJSON(w, 202, map[string]any{"status": "accepted", "command_id": req.CommandID, "command": req.Type, "queued": s.bus != nil})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MessageID  string            `json:"message_id"`
		Channel    string            `json:"channel"`
		SenderID   string            `json:"sender_id"`
		ChatID     string            `json:"chat_id"`
		Content    string            `json:"content"`
		Media      []string          `json:"media,omitempty"`
		SessionKey string            `json:"session_key"`
		Metadata   map[string]string `json:"metadata,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid JSON: " + err.Error(), "code": 400})
		return
	}
	if req.Content == "" {
		writeJSON(w, 400, map[string]any{"error": "message content is required", "code": 400})
		return
	}
	if req.MessageID == "" {
		req.MessageID = fmt.Sprintf("msg-%d", time.Now().UnixNano())
	}
	if req.Channel == "" {
		req.Channel = "edge-api"
	}
	if req.SenderID == "" {
		req.SenderID = "fleet"
	}
	if req.ChatID == "" {
		req.ChatID = req.MessageID
	}
	if req.SessionKey == "" {
		req.SessionKey = "edge:msg:" + req.ChatID
	}
	if req.Metadata == nil {
		req.Metadata = make(map[string]string)
	}
	req.Metadata["kind"] = "message"
	req.Metadata["message_id"] = req.MessageID
	req.Metadata["edge_node_id"] = s.config.NodeID
	if s.bus != nil {
		s.bus.PublishInbound(bus.InboundMessage{
			Channel: req.Channel, SenderID: req.SenderID, ChatID: req.ChatID,
			Content: req.Content, Media: req.Media, SessionKey: req.SessionKey,
			Metadata: req.Metadata,
		})
	}
	log.Printf("[edge] accepted message: id=%s sender=%s node=%s", req.MessageID, req.SenderID, s.config.NodeID)
	writeJSON(w, 202, map[string]any{"status": "accepted", "message_id": req.MessageID, "queued": s.bus != nil})
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
