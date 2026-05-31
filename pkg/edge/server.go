package edge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

// Server is the Edge API HTTP server.
// It exposes endpoints for health checks, sensor data ingestion,
// and command reception from the upstream Fleet Manager.
type Server struct {
	config Config
	mux    *http.ServeMux
	start  time.Time
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

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /api/v1/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/v1/command", s.handleCommand)
}

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	agentName := s.config.AgentName
	if agentName == "" {
		agentName = "picclaw"
	}
	version := s.config.Version
	if version == "" {
		version = "0.1.0"
	}
	buildTime := s.config.BuildTime
	if buildTime == "" {
		buildTime = "dev"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":     "ok",
		"agent":      agentName,
		"agent_name": agentName,
		"node_id":    s.config.NodeID,
		"uptime":     time.Since(s.start).String(),
		"version":    version,
		"build_time": buildTime,
		"go_version": runtime.Version(),
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
