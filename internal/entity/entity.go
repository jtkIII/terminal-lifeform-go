package entity

import (
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/jtkIII/terminal-lifeform-go/internal/data"
)

type Entity struct {
	mu                sync.RWMutex
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Age               float64           `json:"age"`
	Status            data.Status       `json:"status"`
	X, Y              int               `json:"x", "y"`
	Parameters        data.EntityParams `json:"parameters"`
	Health            float64           `json:"health"`
	Energy            float64           `json:"energy"`
	Affinity          string            `json:"affinity"`
	MemorySpan        float64           `json:"memory_span"`
	AdaptationBias    float64           `json:"adaptation_bias"`
	Social            float64           `json:"social"`
	EnvironmentMemory []float64         `json:"environment_memory"`
	Brain             interface{}       `json:"-"` // Placeholder for Brain interface
}

func NewEntity(params *data.EntityParams) *Entity {
	affinity := randomAffinity()
	return &Entity{
		ID:     uuid.New().String()[:8],
		Name:   GenerateName(),
		Age:    rand.Float64() * 21,
		Status: data.StatusAlive,
		X:      0, Y: 0,
		Parameters:        *params,
		Health:            params.InitialHealth,
		Energy:            params.InitialEnergy,
		Affinity:          affinity,
		MemorySpan:        rand.Float64()*20 + 5,
		AdaptationBias:    0.89,
		EnvironmentMemory: make([]float64, 0),
		Social:            (params.Cooperation + params.Curiosity + params.Adaptability) - params.Aggression,
		Brain:             nil, // Will be set by factory
	}
}

func (e *Entity) IsAlive() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Status != data.StatusDead
}

func (e *Entity) UpdateStatus() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Health <= 0 || e.Age >= e.Parameters.MaxAge {
		e.Status = data.StatusDead
		e.Health = 0
		e.Energy = 0
		return
	}

	if e.Energy < 10 {
		e.Status = data.StatusDormant
	} else if e.Energy > 90 && rand.Float64() > 0.9 {
		e.Status = data.StatusExploring
	} else if e.Health >= e.Parameters.ThrivingThresholdHealth && e.Energy >= e.Parameters.ThrivingThresholdEnergy {
		e.Status = data.StatusThriving
		e.Health = min(e.Health+e.Parameters.HealthRecoveryRate*1.33, 90)
	} else if e.Health <= e.Parameters.StrugglingThresholdHealth || e.Energy <= e.Parameters.StrugglingThresholdEnergy {
		e.Status = data.StatusStruggling
		e.Health -= e.Parameters.HealthDecayRate * 1.5
	} else {
		e.Status = data.StatusAlive
	}
}

func randomAffinity() string {
	affinities := []string{"Flux", "Root", "Echo", "Pulse", "Core", "Nixx", "Creed"}
	return affinities[rand.Intn(len(affinities))]
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
