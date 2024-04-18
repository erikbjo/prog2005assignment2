package utils

import (
	"log"
	"os"
)

// DefaultPort Default port for the server
const DefaultPort = "8000"

// TestPort Test port for the server
const TestPort = "8001"

// GetPort Get the port from the environment variable, or use the default port
func GetPort() string {
	// Get the PORT environment variable
	port := os.Getenv("PORT")

	// Use default Port variable if not provided
	if port == "" {
		log.Println("$PORT has not been set. Default: " + DefaultPort)
		port = DefaultPort
	}

	return port
}

// GetTestPort Get the test port from the environment variable, or use the default test port
func GetTestPort() string {
	// Get the PORT environment variable
	port := os.Getenv("TEST_PORT")

	// Use default Port variable if not provided
	if port == "" {
		log.Println("$TEST_PORT has not been set. Default: " + TestPort)
		port = TestPort
	}

	return port
}
