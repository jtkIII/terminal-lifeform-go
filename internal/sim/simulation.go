// internal/sim/simulation.go
package sim

import (
	"context"
	"sync"
	"time"

	"github.com/jtkIII/terminal-lifeform-go/internal/data"
	"github.com/jtkIII/terminal-lifeform-go/internal/entity"
)

// Simulation holds the entire simulation state
type Simulation struct {
	mu sync.RWMutex // Protects all fields below

	// Configuration
	WorldName         string
	TotalEpochs       int
	MemoryWindow      int
	EntropyRate       float64
	MemorySensitivity float64

	// State
	Entities          []*entity.Entity
	EpochCount        int
	PopulationHistory []int
	FeedbackLog       []map[string]interface{}
	Trend             float64
	State             data.EnvironmentState
	MaxEntities       int

	// Control
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// NewSimulation creates a new simulation instance
func NewSimulation(world *data.WorldConfig, initEntities, totalEpochs int) *Simulation {
	ctx, cancel := context.WithCancel(context.Background())

	sim := &Simulation{
		WorldName:         world.Name,
		TotalEpochs:       totalEpochs,
		MemoryWindow:      world.MemoryWindow,
		EntropyRate:       world.EntropyRate,
		MemorySensitivity: world.MemorySensitivity,
		Entities:          make([]*entity.Entity, 0, initEntities),
		PopulationHistory: make([]int, 0, world.MemoryWindow),
		FeedbackLog:       make([]map[string]interface{}, 0, world.MemoryWindow),
		Trend:             0,
		State:             data.StateStable,
		ctx:               ctx,
		cancel:            cancel,
		done:              make(chan struct{}),
	}

	// Initialize entities
	for i := 0; i < initEntities; i++ {
		ent := entity.NewEntity(&data.EntityParams{
			InitialHealth:             50,
			InitialEnergy:             50,
			MaxAge:                    100,
			ThrivingThresholdHealth:   70,
			ThrivingThresholdEnergy:   70,
			StrugglingThresholdHealth: 30,
			StrugglingThresholdEnergy: 30,
			Resilience:                1,
			ForagingEfficiency:        0.1,
			MetabolismRate:            0.1,
			ReproductionChance:        0.05,
			MutationRate:              0.01,
			Aggression:                0.1,
			Cooperation:               0.1,
			HealthRecoveryRate:        1.0,
			HealthDecayRate:           1.0,
			Curiosity:                 0.1,
			Adaptability:              0.1,
		})
		sim.Entities = append(sim.Entities, ent)
	}

	return sim
}

func (s *Simulation) TickOnce() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tickInternal()
}

// tick advances one epoch (called internally by runLoop)
func (s *Simulation) tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tickInternal()
}

// tickInternal contains the actual epoch logic
// Called with lock already held
func (s *Simulation) tickInternal() {
	s.EpochCount++

	// Process each entity
	for _, ent := range s.Entities {
		ent.UpdateStatus()
		if ent.IsAlive() {
			// Process entity logic here
			// ent.Brain.Decide(...) if needed
		}
	}

	// Filter dead entities
	var survivors []*entity.Entity
	for _, ent := range s.Entities {
		if ent.IsAlive() {
			survivors = append(survivors, ent)
		}
	}
	s.Entities = survivors

	// Track population
	pop := len(s.Entities)
	s.PopulationHistory = append(s.PopulationHistory, pop)
	if len(s.PopulationHistory) > s.MemoryWindow {
		s.PopulationHistory = s.PopulationHistory[1:]
	}

	// Update max
	if pop > s.MaxEntities {
		s.MaxEntities = pop
	}

	// Check extinction
	if pop == 0 {
		// Handle extinction scenario
	}

	// Update environment/state logic here
	// handle_enviroment, handle_controllers, etc.
}

// Start begins the simulation loop in a background goroutine
func (s *Simulation) Start() {
	go s.runLoop()
}

// Stop gracefully stops the simulation
func (s *Simulation) Stop() {
	s.cancel()
	<-s.done
}

// runLoop is the main simulation tick loop
func (s *Simulation) runLoop() {
	defer close(s.done)

	ticker := time.NewTicker(1 * time.Second) // Configurable epoch speed
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.tick()
		}
	}
}

// GetStatus returns current simulation state (thread-safe)
func (s *Simulation) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"epoch":          s.EpochCount,
		"population":     len(s.Entities),
		"max_population": s.MaxEntities,
		"state":          s.State,
		"trend":          s.Trend,
		"world":          s.WorldName,
	}
}

// GetEntities returns a copy of current entities (thread-safe)
func (s *Simulation) GetEntities() []*entity.Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return copies to prevent external mutation
	result := make([]*entity.Entity, len(s.Entities))
	copy(result, s.Entities)
	return result
}
