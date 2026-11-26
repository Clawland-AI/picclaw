package edge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server is the Edge API HTTP server.
// It exposes endpoints for health checks, sensor data ingestion,
// and command reception from the upstream Fleet Manager.
type Server struct {
	config Config
	mux    *http.ServeMux
	start  time.Time
	bus    *MessageBus
}

// NewServer creates a new Edge API Server.
func NewServer(cfg Config) *Server {
	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
		start:  time.Now(),
	}
	s.routes()
	return s
}

// SetBus attaches a MessageBus for command/message ingestion.
func (s *Server) SetBus(bus *MessageBus) {
	s.bus = bus
}

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /api/v1/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/v1/command", s.handleCommand)

	// Additional bounty endpoints (non-versioned paths per PicClaw bounty board)
	s.mux.HandleFunc("GET /api/health", s.handleHealthz)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/command", s.handleCommand)
	s.mux.HandleFunc("POST /api/message", s.handleMessage)
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
		"node_id":   s.config.NodeID,
		"node_name": s.config.NodeName,
		"status":    "running",
		"uptime":    time.Since(s.start).String(),
		"cloud":     s.config.Cloud.Endpoint,
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var cmd struct {
		Type    string         `json:"type"`
		Payload map[string]any `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("[edge] received command: type=%s", cmd.Type)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "accepted",
		"command": cmd.Type,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var msg struct {
		Type    string         `json:"type"`
		Payload map[string]any `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("[edge] received message: type=%s", msg.Type)

	// Integrate with MessageBus if available
	if s.bus != nil {
		s.bus.PublishInbound(InboundMessage{
			Type:    msg.Type,
			Payload: msg.Payload,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "accepted",
		"message": msg.Type,
	})
}
