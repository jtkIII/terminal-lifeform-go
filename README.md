# Terminal Lifeform Go
**That's a terrible name, and it's an API now, not term...**

A Go re-write of the inimitable Terminal Lifeform, a Python based sim that ran in the terminal, get it, Terminal Lifeforms live and die in the terminal..  Forget it, needs a new name...

## Simulation Struct (internal/sim/simulation.go)
This is the heart of the system. Key consideration: a background goroutine to run the simulation loop while HTTP handlers can query state.

| -- ### Python Component -- |	-- Go Equivalent --| --	Notes |
| main.py	| cmd/server/main.go	| Entry point, server startup |
| sim.py	| internal/sim/simulation.go	| Core simulation with goroutine loop |
| entity.py |	internal/entity/entity.go |	Entity struct with mutex protection |
| N/A | pkg/api/server.go |	HTTP handlers for API endpoints |
| N/A |	internal/data/types.go |	Shared types/enums |


### Initial Files
terminal-lifeform-go/
├── cmd/
│   └── server/
│       └── main.go          # Entry point (replaces main.py)
├── internal/
│   ├── sim/
│   │   └── simulation.go    # Core Simulation struct (replaces sim.py)
│   ├── entity/
│   │   └── entity.go        # Entity struct (replaces entity.py)
│   ├── world/
│   │   └── world.go         # World/environment config
│   ├── handlers/
│   │   └── api_handlers.go  # HTTP endpoint handlers
│   ├── data/
│   │   └── types.go         # Shared types (Status, Action, etc.)
│   └── utils/
│       └── helpers.go       # Utility functions
├── pkg/
│   └── api/
│       └── server.go        # HTTP server setup
├── go.mod
└── go.sum

