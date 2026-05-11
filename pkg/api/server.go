// pkg/api/server.go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/jtkIII/terminal-lifeform-go/internal/sim"
)

// Server wraps the HTTP server and simulation
type Server struct {
	sim *sim.Simulation
	mux *http.ServeMux
}

// NewServer creates a new API server
func NewServer(s *sim.Simulation) *Server {
	server := &Server{
		sim: s,
		mux: http.NewServeMux(),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/status", s.handleStatus)
	s.mux.HandleFunc("/entities", s.handleEntities)
	s.mux.HandleFunc("/epoch", s.handleEpoch)
	s.mux.HandleFunc("/dump", s.handleDump)
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleStatus returns simulation status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.sim.GetStatus())
}

// handleEntities returns all entities
func (s *Server) handleEntities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.sim.GetEntities())
}

// handleEpoch manually advances one epoch (POST only)
func (s *Server) handleEpoch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Advance simulation by one epoch
	s.sim.TickOnce()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.sim.GetStatus())
}

// handleDump returns full simulation state
func (s *Server) handleDump(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Create a comprehensive dump structure
	dump := map[string]interface{}{
		"status":   s.sim.GetStatus(),
		"entities": s.sim.GetEntities(),
		// Add history, feedback logs, etc.
	}
	json.NewEncoder(w).Encode(dump)
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
