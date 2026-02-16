// Package server provides HTTP health check endpoints for PicoClaw
package server

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

var (
	// StartTime records when the server started
	StartTime time.Time
	// Version holds the application version
	Version = "0.1.0"
	// AgentName holds the agent name
	AgentName = "PicoClaw"
)

func init() {
	StartTime = time.Now()
}

// HealthResponse represents the JSON response from /healthz
type HealthResponse struct {
	Status     string `json:"status"`
	Uptime     string `json:"uptime"`
	Version    string `json:"version"`
	AgentName  string `json:"agent_name"`
	GoVersion  string `json:"go_version"`
	BuildTime  string `json:"build_time,omitempty"`
}

// HealthzHandler handles GET /healthz requests
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := time.Since(StartTime)
	resp := HealthResponse{
		Status:    "ok",
		Uptime:    uptime.String(),
		Version:   Version,
		AgentName: AgentName,
		GoVersion: runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// NewHealthServer creates an HTTP server with the /healthz endpoint
func NewHealthServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthzHandler)
	
	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}