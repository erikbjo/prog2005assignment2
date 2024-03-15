package server

import (
	"assignment-2/server/handlers"
	"assignment-2/server/shared"
	"assignment-2/utils"
	"log"
	"net/http"
)

// Start
/*
Start the server on the port specified in the environment variable PORT. If PORT is not set, the default port 8080 is used.
*/
func Start() {
	// Get the port from the environment variable, or use the default port
	port := utils.GetPort()

	// Set up handler endpoints
	// http.HandleFunc(shared.DefaultPath, handlers.DefaultHandler)
	http.HandleFunc(shared.StatusPath, handlers.StatusHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
