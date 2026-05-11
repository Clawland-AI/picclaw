package edge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
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
func NewServer(cfg Config) *Server {
	return NewServerWithBus(cfg, nil)
}

// NewServerWithBus creates a new Edge API Server wired to the message bus.
func NewServerWithBus(cfg Config, msgBus *bus.MessageBus) *Server {
	cfg = withDefaults(cfg)
	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
		start:  time.Now(),
		bus:    msgBus,
	}
	s.routes()
	return s
}

func withDefaults(cfg Config) Config {
	defaults := DefaultConfig()
	if cfg.Port == 0 {
		cfg.Port = defaults.Port
	}
	if cfg.NodeName == "" {
		cfg.NodeName = defaults.NodeName
	}
	if cfg.Version == "" {
		cfg.Version = defaults.Version
	}
	if cfg.BuildTime == "" {
		cfg.BuildTime = defaults.BuildTime
	}
	if cfg.AgentName == "" {
		cfg.AgentName = defaults.AgentName
	}
	if cfg.Cloud.Endpoint == "" {
		cfg.Cloud.Endpoint = defaults.Cloud.Endpoint
	}
	if cfg.Cloud.HeartbeatInterval == 0 {
		cfg.Cloud.HeartbeatInterval = defaults.Cloud.HeartbeatInterval
	}
	return cfg
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /api/health", s.handleHealthz)
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

// ServeHTTP routes Edge API requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":     "ok",
		"agent_name": s.config.AgentName,
		"node_id":    s.config.NodeID,
		"node_name":  s.config.NodeName,
		"version":    s.config.Version,
		"build_time": s.config.BuildTime,
		"go_version": runtime.Version(),
		"uptime":     time.Since(s.start).String(),
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
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var cmd struct {
		Type       string            `json:"type"`
		Payload    map[string]any    `json:"payload"`
		RequestID  string            `json:"request_id"`
		SessionKey string            `json:"session_key"`
		Metadata   map[string]string `json:"metadata"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	cmd.Type = strings.TrimSpace(cmd.Type)
	if cmd.Type == "" {
		http.Error(w, "command type is required", http.StatusBadRequest)
		return
	}
	if s.bus == nil {
		http.Error(w, "message bus not configured", http.StatusServiceUnavailable)
		return
	}
	if cmd.RequestID == "" {
		cmd.RequestID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if cmd.SessionKey == "" {
		cmd.SessionKey = fmt.Sprintf("edge:%s:%s", s.config.NodeID, cmd.RequestID)
	}
	metadata := cloneMetadata(cmd.Metadata)
	metadata["command_type"] = cmd.Type
	metadata["request_id"] = cmd.RequestID
	if cmd.Payload != nil {
		payload, err := json.Marshal(cmd.Payload)
		if err != nil {
			http.Error(w, "invalid command payload", http.StatusBadRequest)
			return
		}
		metadata["payload"] = string(payload)
	}

	s.bus.PublishInbound(bus.InboundMessage{
		Channel:    "edge-command",
		SenderID:   "fleet-manager",
		ChatID:     s.config.NodeID,
		Content:    cmd.Type,
		SessionKey: cmd.SessionKey,
		Metadata:   metadata,
	})

	log.Printf("[edge] received command: type=%s", cmd.Type)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"status":      "accepted",
		"command":     cmd.Type,
		"request_id":  cmd.RequestID,
		"session_key": cmd.SessionKey,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var msg bus.InboundMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	msg.Content = strings.TrimSpace(msg.Content)
	if msg.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}
	if s.bus == nil {
		http.Error(w, "message bus not configured", http.StatusServiceUnavailable)
		return
	}
	if msg.Channel == "" {
		msg.Channel = "http"
	}
	if msg.SenderID == "" {
		msg.SenderID = "edge-api"
	}
	if msg.ChatID == "" {
		msg.ChatID = s.config.NodeID
	}
	if msg.SessionKey == "" {
		msg.SessionKey = fmt.Sprintf("edge:%s", s.config.NodeID)
	}
	if msg.Metadata == nil {
		msg.Metadata = map[string]string{}
	}

	s.bus.PublishInbound(msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"status":      "accepted",
		"channel":     msg.Channel,
		"session_key": msg.SessionKey,
	})
}

func cloneMetadata(metadata map[string]string) map[string]string {
	cloned := make(map[string]string, len(metadata)+3)
	for key, value := range metadata {
		cloned[key] = value
	}
	return cloned
}
