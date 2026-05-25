// internal/environment/environment.go
package environment

import (
	"sync"
)

// Environment manages the world state
type Environment struct {
	mu sync.RWMutex

	// Factors
	Pollution            float64
	Temperature          float64
	ResourceAvailability float64
	EntropyRate          float64
	MemorySensitivity    float64

	// History
	EnvironmentalMemory []float64
	MemoryWindow        int

	// Calculated
	Trend float64
	State string // stable, boom, bust, collapse
}

// NewEnvironment creates a new environment
func NewEnvironment(config map[string]float64, memoryWindow int) *Environment {
	return &Environment{
		Pollution:            config["pollution"],
		Temperature:          config["temperature"],
		ResourceAvailability: config["resource_availability"],
		EntropyRate:          config["entropy_rate"],
		MemorySensitivity:    config["memory_sensitivity"],
		MemoryWindow:         memoryWindow,
		EnvironmentalMemory:  make([]float64, 0, memoryWindow),
		Trend:                0,
		State:                "stable",
	}
}

// Update processes environmental changes
func (e *Environment) Update(population int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Record current state
	currentState := e.ResourceAvailability - e.Pollution
	e.EnvironmentalMemory = append(e.EnvironmentalMemory, currentState)
	if len(e.EnvironmentalMemory) > e.MemoryWindow {
		e.EnvironmentalMemory = e.EnvironmentalMemory[1:]
	}

	// Calculate trend
	e.Trend = e.calculateTrend()

	// Determine state
	e.State = e.determineState()

	// Natural entropy increase
	e.EntropyRate += 0.001
}

func (e *Environment) calculateTrend() float64 {
	if len(e.EnvironmentalMemory) < 2 {
		return 0
	}

	sum := 0.0
	for i := 1; i < len(e.EnvironmentalMemory); i++ {
		sum += e.EnvironmentalMemory[i] - e.EnvironmentalMemory[i-1]
	}
	return sum / float64(len(e.EnvironmentalMemory)-1)
}

func (e *Environment) determineState() string {
	if e.Trend > 0.5 && e.ResourceAvailability > 70 {
		return "boom"
	} else if e.Trend < -0.5 || e.Pollution > 80 {
		return "collapse"
	} else if e.Trend < 0 && e.ResourceAvailability < 30 {
		return "bust"
	}
	return "stable"
}

// GetFactor returns a specific environmental factor
func (e *Environment) GetFactor(key string) float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	switch key {
	case "pollution":
		return e.Pollution
	case "temperature":
		return e.Temperature
	case "resource_availability":
		return e.ResourceAvailability
	case "entropy":
		return e.EntropyRate
	default:
		return 0
	}
}

// SetFactor sets a specific environmental factor
func (e *Environment) SetFactor(key string, value float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch key {
	case "pollution":
		e.Pollution = clamp(value, 0, 100)
	case "temperature":
		e.Temperature = clamp(value, -50, 100)
	case "resource_availability":
		e.ResourceAvailability = clamp(value, 0, 100)
	}
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
