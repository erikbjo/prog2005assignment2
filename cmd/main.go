// Package main is the entry point for the application, it starts the server.
package main

import (
	"assignment-2/server"
	"github.com/joho/godotenv"
	"log"
)

// main
// Start the server
func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	server.Start()
}
