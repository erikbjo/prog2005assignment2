package server

import (
	"assignment-2/db"
	"assignment-2/server/handlers"
	"assignment-2/server/shared"
	"assignment-2/server/utils"
	"log"
	"net/http"
)

// Start
/*
Start the server on the port specified in the environment variable PORT. If PORT is not set, the default port 8080 is used.
*/
func Start() {
	// Initialization of firebase database
	db.Initialize()

	// Get the port from the environment variable, or use the default port
	port := utils.GetPort()

	// Set up handler endpoints
	http.HandleFunc(shared.StatusPath, handlers.StatusHandler)

	// Serve the web page for any other path
	http.HandleFunc(shared.DefaultPath, handlers.DefaultHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
