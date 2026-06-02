package edge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

// Server is the Edge API HTTP server.
// It exposes endpoints for health checks, sensor data ingestion,
// and command reception from the upstream Fleet Manager.
type Server struct {
	config Config
	mux    *http.ServeMux
	bus    *bus.MessageBus
	start  time.Time

	rateLimiter   *rateLimiter
	requestID     int64
	requestIDMu   sync.Mutex
}

// ServerOption allows functional configuration of the Server.
type ServerOption func(*Server)

// WithMessageBus sets the MessageBus for the server.
func WithMessageBus(mb *bus.MessageBus) ServerOption {
	return func(s *Server) {
		s.bus = mb
	}
}

// WithRateLimit enables per-endpoint rate limiting with the given
// max requests per second and burst size.
func WithRateLimit(rps int, burst int) ServerOption {
	if rps <= 0 {
		rps = 0 // disabled
	}
	return func(s *Server) {
		if rps > 0 {
			s.rateLimiter = newRateLimiter(rps, burst)
		}
	}
}

// NewServer creates a new Edge API Server.
// The variadic bus parameter is kept for backward compatibility with existing callers.
// New code should prefer NewServer(cfg, WithMessageBus(msgBus)).
func NewServer(cfg Config, opts ...any) *Server {
	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
		start:  time.Now(),
	}

	// Process options: accept both *bus.MessageBus (legacy) and ServerOption (new)
	for _, opt := range opts {
		switch v := opt.(type) {
		case *bus.MessageBus:
			s.bus = v
		case ServerOption:
			v(s)
		}
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	// New bounty-required endpoints
	s.mux.HandleFunc("GET /api/health", s.mwChain(s.handleHealthz))
	s.mux.HandleFunc("GET /api/status", s.mwChain(s.handleStatus))
	s.mux.HandleFunc("POST /api/command", s.mwChain(s.handleCommand))
	s.mux.HandleFunc("POST /api/message", s.mwChain(s.handleMessage))

	// Backward-compatible legacy aliases
	s.mux.HandleFunc("GET /healthz", s.mwChain(s.handleHealthz))
	s.mux.HandleFunc("GET /api/v1/status", s.mwChain(s.handleStatus))
	s.mux.HandleFunc("POST /api/v1/command", s.mwChain(s.handleCommand))
}

// ServeHTTP allows tests and embedded callers to use Server as an http.Handler.
// It intercepts OPTIONS preflight requests globally before the mux routing.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight at the top level before mux routing,
	// since Go 1.22+ strict method routing rejects OPTIONS against
	// method-specific routes.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	s.mux.ServeHTTP(w, r)
}

// Start begins listening on the configured port.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
	return http.ListenAndServe(addr, s)
}

// StartWithContext begins listening and returns when the context is cancelled,
// allowing graceful shutdown.
func (s *Server) StartWithContext(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s,
	}

	go func() {
		log.Printf("[edge] API server listening on %s (node: %s)", addr, s.config.NodeID)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[edge] server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}

// Handler returns the underlying http.Handler for embedding into custom servers.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// mwChain applies the server middleware stack to a handler.
func (s *Server) mwChain(next http.HandlerFunc) http.HandlerFunc {
	return s.recoveryMW(s.loggingMW(s.rateLimitMW(next)))
}

// --- Middleware ---

func (s *Server) loggingMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.requestIDMu.Lock()
		s.requestID++
		rid := s.requestID
		s.requestIDMu.Unlock()

		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		r.Header.Set("X-Request-ID", fmt.Sprintf("%d", rid))

		next(lrw, r)

		duration := time.Since(start)
		status := lrw.statusCode
		// Log at INFO level; WARN for 4xx/5xx
		if status >= 500 {
			log.Printf("[edge] %s %s -> %d (%s) rid=%d", r.Method, r.URL.Path, status, duration, rid)
		} else {
			log.Printf("[edge] %s %s -> %d (%s) rid=%d", r.Method, r.URL.Path, status, duration, rid)
		}
	}
}

func (s *Server) recoveryMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[edge] panic recovered: %v", rec)
				writeJSON(w, http.StatusInternalServerError, errorResponse{
					Error: "internal server error",
					Code:  http.StatusInternalServerError,
				})
			}
		}()
		next(w, r)
	}
}

func (s *Server) rateLimitMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.rateLimiter != nil && !s.rateLimiter.Allow() {
			writeJSON(w, http.StatusTooManyRequests, errorResponse{
				Error: "rate limit exceeded",
				Code:  http.StatusTooManyRequests,
			})
			return
		}
		next(w, r)
	}
}

// --- Handlers ---

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	busStatus := "disconnected"
	if s.bus != nil {
		busStatus = "connected"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":        "ok",
		"agent":         "picclaw",
		"node_id":       s.config.NodeID,
		"node_name":     s.config.NodeName,
		"uptime":        time.Since(s.start).String(),
		"version":       "0.1.0",
		"bus_connected": busStatus,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	busInfo := "disconnected"
	if s.bus != nil {
		busInfo = "connected"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"node_id":       s.config.NodeID,
		"node_name":     s.config.NodeName,
		"status":        "running",
		"uptime":        time.Since(s.start).String(),
		"cloud":         s.config.Cloud.Endpoint,
		"bus_connected": busInfo,
		"endpoints": []string{
			"GET /api/health",
			"GET /api/status",
			"POST /api/command",
			"POST /api/message",
			"GET /healthz",
			"GET /api/v1/status",
			"POST /api/v1/command",
		},
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	var req commandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid JSON: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	if req.Type == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "command type is required",
			Code:  http.StatusBadRequest,
		})
		return
	}

	// Normalize defaults
	req.normalize()

	// Build the InboundMessage
	msg := req.toInbound(s.config.NodeID)

	if s.bus != nil {
		s.bus.PublishInbound(msg)
	}

	log.Printf("[edge] accepted command: type=%s id=%s node=%s", req.Type, req.CommandID, s.config.NodeID)

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":     "accepted",
		"command_id": req.CommandID,
		"command":    req.Type,
		"queued":     s.bus != nil,
	})
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var req edgeMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid JSON: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	if req.Content == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "message content is required",
			Code:  http.StatusBadRequest,
		})
		return
	}

	// Normalize defaults
	req.normalize()

	msg := req.toInbound(s.config.NodeID)

	if s.bus != nil {
		s.bus.PublishInbound(msg)
	}

	log.Printf("[edge] accepted message: id=%s sender=%s node=%s", req.MessageID, req.SenderID, s.config.NodeID)

	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":     "accepted",
		"message_id": req.MessageID,
		"queued":     s.bus != nil,
	})
}

// --- Request/Response Types ---

type commandRequest struct {
	CommandID  string            `json:"command_id"`
	Type       string            `json:"type"`
	Payload    map[string]any    `json:"payload"`
	SenderID   string            `json:"sender_id"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

func (c *commandRequest) normalize() {
	if c.CommandID == "" {
		c.CommandID = fmt.Sprintf("cmd-%d", time.Now().UnixNano())
	}
	if c.SenderID == "" {
		c.SenderID = "fleet"
	}
	if c.SessionKey == "" {
		c.SessionKey = "edge:cmd:" + c.CommandID
	}
	if c.Metadata == nil {
		c.Metadata = make(map[string]string)
	}
}

func (c *commandRequest) toInbound(nodeID string) bus.InboundMessage {
	metadata := make(map[string]string, len(c.Metadata)+4)
	for k, v := range c.Metadata {
		metadata[k] = v
	}
	metadata["kind"] = "command"
	metadata["command_id"] = c.CommandID
	metadata["command_type"] = c.Type
	metadata["edge_node_id"] = nodeID
	if c.Payload != nil {
		if payloadBytes, err := json.Marshal(c.Payload); err == nil {
			metadata["payload"] = string(payloadBytes)
		}
	}

	return bus.InboundMessage{
		Channel:    "edge-api",
		SenderID:   c.SenderID,
		ChatID:     c.CommandID,
		Content:    c.Type,
		SessionKey: c.SessionKey,
		Metadata:   metadata,
	}
}

type edgeMessageRequest struct {
	MessageID  string            `json:"message_id"`
	Channel    string            `json:"channel"`
	SenderID   string            `json:"sender_id"`
	ChatID     string            `json:"chat_id"`
	Content    string            `json:"content"`
	Media      []string          `json:"media,omitempty"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

func (m *edgeMessageRequest) normalize() {
	if m.MessageID == "" {
		m.MessageID = fmt.Sprintf("msg-%d", time.Now().UnixNano())
	}
	if m.Channel == "" {
		m.Channel = "edge-api"
	}
	if m.SenderID == "" {
		m.SenderID = "fleet"
	}
	if m.ChatID == "" {
		m.ChatID = m.MessageID
	}
	if m.SessionKey == "" {
		m.SessionKey = "edge:msg:" + m.ChatID
	}
	if m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
}

func (m *edgeMessageRequest) toInbound(nodeID string) bus.InboundMessage {
	metadata := make(map[string]string, len(m.Metadata)+3)
	for k, v := range m.Metadata {
		metadata[k] = v
	}
	metadata["kind"] = "message"
	metadata["message_id"] = m.MessageID
	metadata["edge_node_id"] = nodeID

	return bus.InboundMessage{
		Channel:    m.Channel,
		SenderID:   m.SenderID,
		ChatID:     m.ChatID,
		Content:    m.Content,
		Media:      m.Media,
		SessionKey: m.SessionKey,
		Metadata:   metadata,
	}
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// loggingResponseWriter wraps http.ResponseWriter to capture status codes.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// --- Rate Limiter ---

type rateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	interval time.Duration
	last     time.Time
}

func newRateLimiter(rps int, burst int) *rateLimiter {
	rl := &rateLimiter{
		tokens:   float64(burst),
		max:      float64(burst),
		interval: time.Second / time.Duration(rps),
		last:     time.Now(),
	}
	// If burst < 1, set to at least rps
	if burst < 1 {
		rl.max = float64(rps)
		rl.tokens = float64(rps)
	}
	return rl
}

func (rl *rateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.last)
	rl.last = now

	// Refill tokens based on elapsed time
	rl.tokens += elapsed.Seconds() * (rl.max / rl.interval.Seconds())
	if rl.tokens > rl.max {
		rl.tokens = rl.max
	}

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}
