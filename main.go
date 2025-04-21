package main

import (
	"log"

	"personatrip/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
