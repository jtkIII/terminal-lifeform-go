package brain

import (
	"math/rand"

	"github.com/jtkIII/terminal-lifeform-go/internal/data"
	"github.com/jtkIII/terminal-lifeform-go/internal/entity"
)

type SimulationRef interface {
	GetEnvFactor(key string) float64
	GetEntities() []*entity.Entity
	GetEntropy() float64
}

type Brain interface {
	Decide(ent *entity.Entity, sim SimulationRef) data.Action
}

type EntMind struct {
	Curiosity float64
	Fear      float64
	Goal      string
}

func NewEntMind() *EntMind {
	return &EntMind{Curiosity: 0.5, Fear: 0.2, Goal: "idle"}
}

func getEntropy(affinity string) float64 {
	entropyMap := map[string]float64{
		"Flux": 0.2, "Root": 0.05, "Echo": 0.1, "Pulse": 0.15,
		"Core": 0.1, "Nixx": 0.05, "Creed": 0.1,
	}
	if val, ok := entropyMap[affinity]; ok {
		return val
	}
	return 0.1
}

type FluxBrain struct{}

func (b *FluxBrain) Decide(ent *entity.Entity, sim SimulationRef) data.Action {
	r := rand.Float64() + getEntropy(ent.Affinity) + sim.GetEntropy()/3
	if r < 0.6 {
		return data.ActionExplore
	}
	if r < 0.8 {
		return data.ActionInspect
	}
	return data.ActionWait
}

type RootBrain struct{}

func (b *RootBrain) Decide(ent *entity.Entity, sim SimulationRef) data.Action {
	if sim.GetEnvFactor("resource_availability") < 40 {
		return data.ActionGather
	}
	if ent.Energy < 30 {
		return data.ActionRest
	}
	if ent.Health < 50 {
		return data.ActionHeal
	}
	return data.ActionWait
}

type EchoBrain struct{}

func (b *EchoBrain) Decide(ent *entity.Entity, sim SimulationRef) data.Action {
	hasStruggling := false
	for _, e := range sim.GetEntities() {
		if e.ID != ent.ID && e.Status == data.StatusStruggling {
			hasStruggling = true
			break
		}
	}
	if hasStruggling {
		return data.ActionCommunicate
	}
	if ent.Energy > 70 && ent.Health > 70 {
		return data.ActionSocialize
	}
	return data.ActionWait
}

type PulseBrain struct{}

func (b *PulseBrain) Decide(ent *entity.Entity, sim SimulationRef) data.Action {
	if ent.Energy < 40 {
		return data.ActionGather
	}
	if ent.Health < 60 {
		return data.ActionHeal
	}
	if rand.Float64() < 0.3 {
		return data.ActionMove
	}
	return data.ActionWait
}

type CoreBrain struct{}

func (b *CoreBrain) Decide(ent *entity.Entity, sim SimulationRef) data.Action {
	if ent.Health < 50 {
		return data.ActionHeal
	}
	if ent.Energy < 50 {
		return data.ActionGather
	}
	if rand.Float64() < 0.2 {
		return data.ActionExplore
	}
	return data.ActionWait
}

type NixxBrain struct{}

func (b *NixxBrain) Decide(ent *entity.Entity, sim SimulationRef) data.Action {
	if ent.Health < 40 {
		return data.ActionHide
	}
	if sim.GetEnvFactor("pollution") > 70 {
		return data.ActionRelocate
	}
	return data.ActionWait
}

type CreedBrain struct{}

func (b *CreedBrain) Decide(ent *entity.Entity, sim SimulationRef) data.Action {
	hasWeaker := false
	for _, e := range sim.GetEntities() {
		if e.ID != ent.ID && e.Health < ent.Health {
			hasWeaker = true
			break
		}
	}
	if hasWeaker {
		return data.ActionAttack
	}
	if ent.Energy > 80 {
		return data.ActionPatrol
	}
	return data.ActionWait
}

var AffinityBrains = map[string]func() Brain{
	"Flux":  func() Brain { return &FluxBrain{} },
	"Root":  func() Brain { return &RootBrain{} },
	"Echo":  func() Brain { return &EchoBrain{} },
	"Pulse": func() Brain { return &PulseBrain{} },
	"Core":  func() Brain { return &CoreBrain{} },
	"Nixx":  func() Brain { return &NixxBrain{} },
	"Creed": func() Brain { return &CreedBrain{} },
}

func CreateBrain(affinity string) Brain {
	if constructor, ok := AffinityBrains[affinity]; ok {
		return constructor()
	}
	return &CoreBrain{}
}
