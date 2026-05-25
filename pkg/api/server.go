package api

import (
	"encoding/json"
	"net/http"

	"github.com/jtkIII/terminal-lifeform-go/internal/sim"
)

type Server struct {
	sim *sim.Simulation
	mux *http.ServeMux
}

func NewServer(s *sim.Simulation) *Server {
	server := &Server{sim: s, mux: http.NewServeMux()}
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

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.sim.GetStatus())
}

func (s *Server) handleEntities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.sim.GetEntitiesPublic())
}

func (s *Server) handleEpoch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.sim.TickOnce()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.sim.GetStatus())
}

func (s *Server) handleDump(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dump := map[string]interface{}{
		"status":   s.sim.GetStatus(),
		"entities": s.sim.GetEntities(),
	}
	json.NewEncoder(w).Encode(dump)
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
