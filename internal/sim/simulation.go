package sim

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/jtkIII/terminal-lifeform-go/internal/brain"
	"github.com/jtkIII/terminal-lifeform-go/internal/data"
	"github.com/jtkIII/terminal-lifeform-go/internal/entity"
	"github.com/jtkIII/terminal-lifeform-go/internal/handlers"
)

// type EventType string

// const (
// 	EventBirth     EventType = "birth"
// 	EventDeath     EventType = "death"
// 	EventWar       EventType = "war"
// 	EventBabyBoom  EventType = "baby_boom"
// 	EventActOfGod  EventType = "act_of_god"
// 	EventReproduce EventType = "reproduce"
// )

type SimulationEvent struct {
	Type      data.EventType `json:"type"`
	Epoch     int            `json:"epoch"`
	EntityID  string         `json:"entity_id,omitempty"`
	Details   string         `json:"details,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

type Simulation struct {
	mu                sync.RWMutex
	WorldName         string
	TotalEpochs       int
	MemoryWindow      int
	EntropyRate       float64
	MemorySensitivity float64
	Entities          []*entity.Entity
	EpochCount        int
	PopulationHistory []int
	FeedbackLog       []map[string]interface{}
	Trend             float64
	State             data.EnvironmentState
	MaxEntities       int
	EnvFactors        map[string]float64
	ctx               context.Context
	cancel            context.CancelFunc
	done              chan struct{}
	Handlers          *handlers.Handlers
	Events            []SimulationEvent `json:"-"` // Not serialized in dump, but accessible via API

}

func NewSimulation(world *data.WorldConfig, initEntities, totalEpochs int) *Simulation {
	ctx, cancel := context.WithCancel(context.Background())

	// 1. Create the base simulation struct FIRST
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
		EnvFactors: map[string]float64{
			"resource_availability": world.ResourceAvailability,
			"pollution":             world.Pollution,
			"temperature":           world.Temperature,
		},
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
		// NOTE: Do NOT set Handlers here yet!
	}

	// 2. Initialize entities
	for range initEntities {
		params := &data.EntityParams{
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
		}
		ent := entity.NewEntity(params)

		// Assign the brain
		ent.Brain = brain.CreateBrain(ent.Affinity)

		sim.Entities = append(sim.Entities, ent)
	}

	// 3. NOW assign the Handlers (sim exists at this point!)
	sim.Handlers = handlers.NewHandlers(sim)

	return sim
}

func (s *Simulation) IncrementEpoch() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EpochCount++
}

func (s *Simulation) SetEnvFactor(key string, value float64) {
	if s.EnvFactors == nil {
		s.EnvFactors = make(map[string]float64)
	}
	s.EnvFactors[key] = value
}

func (s *Simulation) Start() { go s.runLoop() }
func (s *Simulation) Stop()  { s.cancel(); <-s.done }

func (s *Simulation) runLoop() {
	defer close(s.done)
	ticker := time.NewTicker(1 * time.Second)
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

func (s *Simulation) tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.EpochCount++

	s.Handlers.HandleEnvironment()
	s.Handlers.HandleInteractions(s.Entities)

	s.Entities = s.Handlers.HandleReproduction(s.Entities)
	s.Entities = s.Handlers.HandleBabyBoom(s.Entities)

	// 	process individual entity status updates
	for _, ent := range s.Entities {
		ent.UpdateStatus()
	}

	// handle special events (war/god)
	if rand.Float64() < 0.05 { // 5% chance per epoch
		s.Handlers.HandleActOfWar(s.Entities)

	}
	if rand.Float64() < 0.02 { // 2% chance per epoch
		s.Handlers.HandleActOfGod(s.Entities)
	}

	// Filter dead entities
	var survivors []*entity.Entity
	for _, ent := range s.Entities {
		if ent.IsAlive() {
			survivors = append(survivors, ent)
		} else {
			s.LogEvent(data.EventDeath, ent.ID, fmt.Sprintf("%s died at age %.1f", ent.Name, ent.Age))
		}
	}
	s.Entities = survivors

	// Track population
	pop := len(s.Entities)
	s.PopulationHistory = append(s.PopulationHistory, pop)
	if len(s.PopulationHistory) > s.MemoryWindow {
		s.PopulationHistory = s.PopulationHistory[1:]
	}

	if pop > s.MaxEntities {
		s.MaxEntities = pop
	}
}

func (s *Simulation) TickOnce() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tick()
}

// Returns a specific environment factor (NO LOCK - caller must hold lock)
func (s *Simulation) GetEnvFactor(key string) float64 {
	if s.EnvFactors == nil {
		return 0.0
	}
	if val, ok := s.EnvFactors[key]; ok {
		return val
	}
	return 0.0
}

// Returns a copy of entities (NO LOCK - caller must hold lock)
func (s *Simulation) GetEntities() []*entity.Entity {
	if s.Entities == nil {
		return []*entity.Entity{}
	}
	result := make([]*entity.Entity, len(s.Entities))
	copy(result, s.Entities)
	return result
}

// Returns the current entropy rate (NO LOCK - caller must hold lock)
func (s *Simulation) GetEntropy() float64 {
	return s.EntropyRate
}

// GetStatus DOES lock because it's called by API (external caller)
func (s *Simulation) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// We need to make a copy of entities for the status if we return them
	// But for status, we usually just return counts.
	return map[string]interface{}{
		"epoch":          s.EpochCount,
		"population":     len(s.Entities),
		"max_population": s.MaxEntities,
		"state":          s.State,
		"trend":          s.Trend,
		"world":          s.WorldName,
	}
}

// The public API version that locks
func (s *Simulation) GetEntitiesPublic() []*entity.Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*entity.Entity, len(s.Entities))
	copy(result, s.Entities)
	return result
}

// LogEvent logs a simulation event
func (s *Simulation) LogEvent(eventType data.EventType, entityID string, details string) {
	s.Events = append(s.Events, SimulationEvent{
		Type:      eventType,
		Epoch:     s.EpochCount,
		EntityID:  entityID,
		Details:   details,
		Timestamp: time.Now(),
	})

	if len(s.Events) > 1000 {
		s.Events = s.Events[len(s.Events)-1000:]
	}
}
