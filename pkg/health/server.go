package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// Server provides HTTP health check endpoints
type Server struct {
	port      int
	startTime time.Time
	version   string
	buildTime string
	agentName string
}

// Config holds health server configuration
type Config struct {
	Port      int
	Version   string
	BuildTime string
	AgentName string
}

// NewServer creates a new health check server
func NewServer(cfg Config) *Server {
	return &Server{
		port:      cfg.Port,
		startTime: time.Now(),
		version:   cfg.Version,
		buildTime: cfg.BuildTime,
		agentName: cfg.AgentName,
	}
}

// Start begins listening on the configured port
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	
	addr := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(addr, mux)
}

// handleHealthz returns health status
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := time.Since(s.startTime)
	
	response := map[string]interface{}{
		"status":     "ok",
		"uptime":     uptime.String(),
		"version":    s.version,
		"agent_name": s.agentName,
		"go_version": runtime.Version(),
		"build_time": s.buildTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
