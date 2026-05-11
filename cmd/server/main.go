// cmd/server/main.go
package main

import (
	"log"
	"os"

	"github.com/jtkIII/terminal-lifeform-go/internal/data"
	"github.com/jtkIII/terminal-lifeform-go/internal/sim"
	"github.com/jtkIII/terminal-lifeform-go/pkg/api"
)

func main() {
	// Load world config (could come from env, file, or CLI arg)
	world := &data.WorldConfig{
		Name:                 "default",
		MemoryWindow:         20,
		MemorySensitivity:    1.1,
		EntropyRate:          0.05,
		Temperature:          20,
		Pollution:            0,
		ResourceAvailability: 1.0,
		MutationRate:         0.01,
	}

	// Create simulation
	s := sim.NewSimulation(world, 100, 1000)
	s.Start() // Begin background simulation loop

	// Create API server
	server := api.NewServer(s)

	// Listen for shutdown signals
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Terminal Lifeform API on :%s", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
