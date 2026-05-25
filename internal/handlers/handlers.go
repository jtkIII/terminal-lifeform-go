// internal/handlers/handlers.go
package handlers

import (
	"math/rand"
	"sync"

	"github.com/jtkIII/terminal-lifeform-go/internal/entity"
)

// SimulationRef provides the handlers with access to simulation state
type SimulationRef interface {
	GetEntities() []*entity.Entity
	GetEnvFactor(key string) float64
	SetEnvFactor(key string, value float64)
	GetEntropy() float64
	// IncrementEpoch()
}

// Handlers manages all entity interactions and events
type Handlers struct {
	mu  sync.Mutex
	sim SimulationRef
}

// NewHandlers creates a new handlers instance
func NewHandlers(sim SimulationRef) *Handlers {
	return &Handlers{sim: sim}
}

// HandleInteractions processes all entity-to-entity interactions
func (h *Handlers) HandleInteractions(entities []*entity.Entity) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			e1 := entities[i]
			e2 := entities[j]

			if !e1.IsAlive() || !e2.IsAlive() {
				continue
			}

			h.processInteraction(e1, e2)
		}
	}
}

func (h *Handlers) processInteraction(e1, e2 *entity.Entity) {
	// Check proximity (simplified - in real impl, check X,Y coordinates)
	distance := h.calculateDistance(e1, e2)
	if distance > 5 {
		return // Too far apart to interact
	}

	// Affinity-based interaction
	if e1.Affinity == e2.Affinity {
		h.handleSameAffinity(e1, e2)
	} else {
		h.handleDifferentAffinity(e1, e2)
	}
}

func (h *Handlers) calculateDistance(e1, e2 *entity.Entity) float64 {
	dx := float64(e1.X - e2.X)
	dy := float64(e1.Y - e2.Y)
	return (dx*dx + dy*dy) // Simplified distance
}

func (h *Handlers) handleSameAffinity(e1, e2 *entity.Entity) {
	// Same affinity entities cooperate
	if rand.Float64() < 0.1 {
		// Small chance to heal each other
		e1.Health = min(e1.Health+5, 100)
		e2.Health = min(e2.Health+5, 100)
	}
}

func (h *Handlers) handleDifferentAffinity(e1, e2 *entity.Entity) {
	// Different affinities may compete or attack
	if e1.Affinity == "Creed" && e2.Affinity != "Creed" {
		// Creed attacks others
		if rand.Float64() < e1.Parameters.Aggression {
			e2.Health -= 10
			e1.Energy += 5
		}
	} else if e1.Affinity == "Root" && e2.Affinity == "Flux" {
		// Root protects against Flux exploration
		e1.Energy -= 2
	}
}

// HandleReproduction processes entity reproduction
func (h *Handlers) HandleReproduction(entities []*entity.Entity) []*entity.Entity {
	h.mu.Lock()
	defer h.mu.Unlock()

	var newEntities []*entity.Entity
	resourceAvail := h.sim.GetEnvFactor("resource_availability")

	for _, ent := range entities {
		if !ent.IsAlive() {
			continue
		}

		// Check reproduction conditions
		if ent.Energy > 80 && ent.Health > 70 && ent.Age < ent.Parameters.MaxAge*0.7 {
			if rand.Float64() < ent.Parameters.ReproductionChance*resourceAvail {
				child := h.createChild(ent)
				newEntities = append(newEntities, child)
			}
		}
	}

	// Add new entities to simulation
	for _, child := range newEntities {
		entities = append(entities, child)
	}

	return entities
}

func (h *Handlers) createChild(parent *entity.Entity) *entity.Entity {
	// Create child with slight mutations
	params := parent.Parameters
	params.MutationRate += (rand.Float64() - 0.5) * 0.1 // Small mutation

	child := entity.NewEntity(&params)
	child.Age = 0
	child.Health = parent.Health * 0.8
	child.Energy = parent.Energy * 0.5
	parent.Energy *= 0.5 // Parent loses energy

	// Inherit affinity with possible mutation
	if rand.Float64() < params.MutationRate {
		child.Affinity = randomMutatedAffinity(parent.Affinity)
	}

	return child
}

func randomMutatedAffinity(parentAffinity string) string {
	affinities := []string{"Flux", "Root", "Echo", "Pulse", "Core", "Nixx", "Creed"}
	// 30% chance to mutate to different affinity
	if rand.Float64() < 0.3 {
		for _, a := range affinities {
			if a != parentAffinity {
				return a
			}
		}
	}
	return parentAffinity
}

// HandleBabyBoom processes special reproduction events
func (h *Handlers) HandleBabyBoom(entities []*entity.Entity) []*entity.Entity {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Baby boom happens when conditions are ideal
	resourceAvail := h.sim.GetEnvFactor("resource_availability")
	if resourceAvail > 80 && rand.Float64() < 0.05 {
		// 10% of healthy entities reproduce
		for _, ent := range entities {
			if ent.IsAlive() && ent.Energy > 60 && ent.Health > 60 {
				if rand.Float64() < 0.1 {
					child := h.createChild(ent)
					entities = append(entities, child)
				}
			}
		}
	}

	return entities
}

// HandleActOfWar processes aggressive events
func (h *Handlers) HandleActOfWar(entities []*entity.Entity) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Find all Creed entities
	var creeds []*entity.Entity
	for _, ent := range entities {
		if ent.Affinity == "Creed" && ent.IsAlive() {
			creeds = append(creeds, ent)
		}
	}

	// Creeds attack weaker entities
	for _, creed := range creeds {
		for _, target := range entities {
			if target.ID != creed.ID && target.IsAlive() && target.Health < creed.Health {
				if rand.Float64() < 0.3 {
					target.Health -= 20
					creed.Energy += 10
				}
			}
		}
	}
}

// HandleActOfGod processes environmental disaster events
func (h *Handlers) HandleActOfGod(entities []*entity.Entity) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Random chance of environmental event
	if rand.Float64() < 0.02 {
		eventType := rand.Intn(4)
		switch eventType {
		case 0: // Meteor strike - random damage
			for _, ent := range entities {
				if rand.Float64() < 0.1 {
					ent.Health -= 30
				}
			}
		case 1: // Resource surge
			h.sim.SetEnvFactor("resource_availability", h.sim.GetEnvFactor("resource_availability")+20)
		case 2: // Pollution spike
			h.sim.SetEnvFactor("pollution", h.sim.GetEnvFactor("pollution")+15)
		case 3: // Temperature shift
			h.sim.SetEnvFactor("temperature", h.sim.GetEnvFactor("temperature")+5)
		}
	}
}

// HandleEnvironment processes environmental changes
func (h *Handlers) HandleEnvironment() {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Natural entropy increase
	currentPollution := h.sim.GetEnvFactor("pollution")
	currentResource := h.sim.GetEnvFactor("resource_availability")

	// Pollution naturally increases
	h.sim.SetEnvFactor("pollution", currentPollution+0.5)

	// Resources deplete with population
	population := len(h.sim.GetEntities())
	if population > 50 {
		h.sim.SetEnvFactor("resource_availability", currentResource-0.3)
	}

	// Clamp values
	if h.sim.GetEnvFactor("pollution") > 100 {
		h.sim.SetEnvFactor("pollution", 100)
	}
	if h.sim.GetEnvFactor("resource_availability") < 0 {
		h.sim.SetEnvFactor("resource_availability", 0)
	}
}

// Helper function
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
