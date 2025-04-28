package main

import (
	"personatrip/cmd"
	"personatrip/internal/utils/logger"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
