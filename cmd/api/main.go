// Package main is the entry point for the application, it starts the server.
package main

import (
	"assignment-2/internal/config"
	`assignment-2/internal/http/server`
	"log"
)

func init() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}
	log.Println("Configuration loaded successfully")
}

// main
// Start the server
func main() {
	// Start the server
	server.Start()
}
