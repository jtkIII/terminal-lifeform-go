// internal/entity/entity.go
package entity

import (
	"math/rand"
	"sync"
	"time"

	"github.com/jtkIII/terminal-lifeform-go/internal/data"

	"github.com/google/uuid"
)

// Entity represents an individual lifeform in the simulation
type Entity struct {
	mu sync.RWMutex // Protects all fields below

	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Age    float64     `json:"age"`
	Status data.Status `json:"status"`
	X, Y   int         `json:"x", "y"`

	// Parameters
	Parameters data.EntityParams `json:"parameters"`

	// State
	Health         float64 `json:"health"`
	Energy         float64 `json:"energy"`
	Affinity       string  `json:"affinity"`
	MemorySpan     float64 `json:"memory_span"`
	AdaptationBias float64 `json:"adaptation_bias"`
	Social         float64 `json:"social"`

	// History
	EnvironmentMemory []float64 `json:"environment_memory"`

	// Brain (placeholder - we'll define brain interface separately)
	Brain Brain `json:"-"` // Not serialized
}

// Brain interface for entity decision-making
type Brain interface {
	Decide(e *Entity, sim interface{}) data.Action
}

// NewEntity creates a new entity with randomized attributes
func NewEntity(params *data.EntityParams) *Entity {
	now := time.Now()
	_ = now // Placeholder for name generation

	return &Entity{
		ID:                uuid.New().String()[:8],
		Name:              generateName(), // We'll implement this
		Age:               rand.Float64() * 21,
		Status:            data.StatusAlive,
		X:                 0, // Will be set by world
		Y:                 0,
		Parameters:        *params,
		Health:            params.InitialHealth,
		Energy:            params.InitialEnergy,
		Affinity:          randomAffinity(),
		MemorySpan:        rand.Float64()*20 + 5,
		AdaptationBias:    0.89,
		EnvironmentMemory: make([]float64, 0),
		Social:            (params.Cooperation + params.Curiosity + params.Adaptability) - params.Aggression,
	}
}

// IsAlive checks if the entity is still living
func (e *Entity) IsAlive() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Status != data.StatusDead
}

// UpdateStatus updates entity status based on health/energy
func (e *Entity) UpdateStatus() {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Death check
	if e.Health <= 0 || e.Age >= e.Parameters.MaxAge {
		e.Status = data.StatusDead
		e.Health = 0
		e.Energy = 0
		return
	}

	// Status determination
	if e.Energy < 10 {
		e.Status = data.StatusDormant
	} else if e.Energy > 90 && rand.Float64() > 0.9 {
		e.Status = data.StatusExploring
	} else if e.Health >= e.Parameters.ThrivingThresholdHealth &&
		e.Energy >= e.Parameters.ThrivingThresholdEnergy {
		e.Status = data.StatusThriving
		e.Health = min(e.Health+e.Parameters.HealthRecoveryRate*1.33, 90)
	} else if e.Health <= e.Parameters.StrugglingThresholdHealth ||
		e.Energy <= e.Parameters.StrugglingThresholdEnergy {
		e.Status = data.StatusStruggling
		e.Health -= e.Parameters.HealthDecayRate * 1.5
	} else {
		e.Status = data.StatusAlive
	}

	e.handleAffinity()
	e.shiftAffinity()
}

// Helper functions (implement these separately)
// func generateName() string        { /* ... */ return "Ent" }
// func randomAffinity() string      { /* ... */ return "default" }
func (e *Entity) handleAffinity() { /* ... */ }
func (e *Entity) shiftAffinity()  { /* ... */ }

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Add these helper functions at the bottom of internal/entity/entity.go

var names = []string{
	"Aria", "Bex", "Cael", "Dax", "Evo", "Flux", "Grit", "Halo",
	"Iris", "Jinx", "Kai", "Luna", "Mox", "Nyx", "Orion", "Pax",
	"Quin", "Rex", "Sage", "Tess", "Umbra", "Vex", "Wren", "Zed",
}

var affinities = []string{"Flux", "Root", "Echo", "Pulse", "Core", "Nixx", "Creed"}

func generateName() string {
	return names[rand.Intn(len(names))]
}

func randomAffinity() string {
	return affinities[rand.Intn(len(affinities))]
}
