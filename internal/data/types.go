package data

type Status string

const (
	StatusAlive      Status = "alive"
	StatusDead       Status = "dead"
	StatusDormant    Status = "dormant"
	StatusExploring  Status = "exploring"
	StatusThriving   Status = "thriving"
	StatusStruggling Status = "struggling"
)

type Action string

const (
	ActionExplore     Action = "explore"
	ActionMove        Action = "move"
	ActionInspect     Action = "inspect"
	ActionWait        Action = "wait"
	ActionGather      Action = "gather"
	ActionCommunicate Action = "communicate"
	ActionHeal        Action = "heal"
	ActionRest        Action = "rest"
	ActionHide        Action = "hide"
	ActionAttack      Action = "attack"
	ActionSocialize   Action = "socialize"
	ActionPatrol      Action = "patrol"
	ActionRelocate    Action = "relocate"
)

type EnvironmentState string

const (
	StateStable   EnvironmentState = "stable"
	StateBoom     EnvironmentState = "boom"
	StateBust     EnvironmentState = "bust"
	StateCollapse EnvironmentState = "collapse"
)

type EntityParams struct {
	InitialHealth             float64 `json:"initial_health"`
	InitialEnergy             float64 `json:"initial_energy"`
	MaxAge                    float64 `json:"max_age"`
	ThrivingThresholdHealth   float64 `json:"thriving_threshold_health"`
	ThrivingThresholdEnergy   float64 `json:"thriving_threshold_energy"`
	StrugglingThresholdHealth float64 `json:"struggling_threshold_health"`
	StrugglingThresholdEnergy float64 `json:"struggling_threshold_energy"`
	Resilience                float64 `json:"resilience"`
	ForagingEfficiency        float64 `json:"foraging_efficiency"`
	MetabolismRate            float64 `json:"metabolism_rate"`
	ReproductionChance        float64 `json:"reproduction_chance"`
	MutationRate              float64 `json:"mutation_rate"`
	Aggression                float64 `json:"aggression"`
	Cooperation               float64 `json:"cooperation"`
	HealthRecoveryRate        float64 `json:"health_recovery_rate"`
	HealthDecayRate           float64 `json:"health_decay_rate"`
	Curiosity                 float64 `json:"curiosity"`
	Adaptability              float64 `json:"adaptability"`
}

type WorldConfig struct {
	Name                 string  `json:"name"`
	MemoryWindow         int     `json:"memory_window"`
	MemorySensitivity    float64 `json:"memory_sensitivity"`
	EntropyRate          float64 `json:"entropy_rate"`
	Temperature          float64 `json:"temperature"`
	Pollution            float64 `json:"pollution"`
	ResourceAvailability float64 `json:"resource_availability"`
	MutationRate         float64 `json:"mutation_rate"`
}
