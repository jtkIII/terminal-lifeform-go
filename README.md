# Terminal Lifeform Go
**That's a terrible name, and it's an API now, not term...**

A Go re-write of the inimitable Terminal Lifeform, a Python based sim that ran in the terminal, get it, Terminal Lifeforms live and die in the terminal..  Forget it, needs a new name...

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

