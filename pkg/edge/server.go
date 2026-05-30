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
	msgBus := bus.NewMessageBus()
	if len(messageBus) > 0 && messageBus[0] != nil {
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
	s.mux.HandleFunc("GET /api/health", s.handleHealthz)
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("GET /api/v1/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/command", s.handleCommand)
	s.mux.HandleFunc("POST /api/v1/command", s.handleCommand)
	s.mux.HandleFunc("POST /api/message", s.handleMessage)
}

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"agent":   "picclaw",
		"node_id": s.config.NodeID,
		"uptime":  time.Since(s.start).String(),
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"node_id":     s.config.NodeID,
		"node_name":   s.config.NodeName,
		"status":      "running",
		"uptime":      time.Since(s.start).String(),
		"cloud":       s.config.Cloud.Endpoint,
		"message_bus": "connected",
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID         string         `json:"id"`
		Type       string         `json:"type"`
		Payload    map[string]any `json:"payload"`
		SenderID   string         `json:"sender_id"`
		ChatID     string         `json:"chat_id"`
		SessionKey string         `json:"session_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		http.Error(w, "missing command type", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		req.ID = fmt.Sprintf("cmd-%d", time.Now().UnixNano())
	}
	if req.SenderID == "" {
		req.SenderID = "fleet"
	}
	if req.ChatID == "" {
		req.ChatID = req.ID
	}
	if req.SessionKey == "" {
		req.SessionKey = "edge-command:" + req.ID
	}

	content, err := json.Marshal(map[string]any{
		"id":      req.ID,
		"type":    req.Type,
		"payload": req.Payload,
	})
	if err != nil {
		http.Error(w, "invalid command payload", http.StatusBadRequest)
		return
	}

	s.bus.PublishInbound(bus.InboundMessage{
		Channel:    "edge-api",
		SenderID:   req.SenderID,
		ChatID:     req.ChatID,
		Content:    string(content),
		SessionKey: req.SessionKey,
		Metadata: map[string]string{
			"kind":         "command",
			"command_id":   req.ID,
			"command_type": req.Type,
			"source":       "edge-http",
		},
	})

	log.Printf("[edge] queued command: id=%s type=%s", req.ID, req.Type)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"status":     "accepted",
		"command_id": req.ID,
		"command":    req.Type,
		"queued":     true,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Channel    string            `json:"channel"`
		SenderID   string            `json:"sender_id"`
		ChatID     string            `json:"chat_id"`
		Content    string            `json:"content"`
		Media      []string          `json:"media,omitempty"`
		SessionKey string            `json:"session_key"`
		Metadata   map[string]string `json:"metadata,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "missing message content", http.StatusBadRequest)
		return
	}

	if req.Channel == "" {
		req.Channel = "edge-api"
	}
	if req.SenderID == "" {
		req.SenderID = "fleet"
	}
	if req.ChatID == "" {
		req.ChatID = s.config.NodeID
	}
	if req.SessionKey == "" {
		req.SessionKey = "edge-message:" + req.ChatID
	}
	if req.Metadata == nil {
		req.Metadata = make(map[string]string)
	}
	req.Metadata["kind"] = "message"
	req.Metadata["source"] = "edge-http"

	s.bus.PublishInbound(bus.InboundMessage{
		Channel:    req.Channel,
		SenderID:   req.SenderID,
		ChatID:     req.ChatID,
		Content:    req.Content,
		Media:      req.Media,
		SessionKey: req.SessionKey,
		Metadata:   req.Metadata,
	})

	log.Printf("[edge] queued message: channel=%s sender=%s", req.Channel, req.SenderID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"status":      "accepted",
		"channel":     req.Channel,
		"session_key": req.SessionKey,
		"queued":      true,
	})
}
