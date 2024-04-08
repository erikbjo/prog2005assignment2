package server

import (
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
	// Get the port from the environment variable, or use the default port
	port := utils.GetPort()

	// Using mux to handle /'s and parameters
	mux := http.NewServeMux()

	// Set up handler endpoints, with and without trailing slash
	// Status
	mux.HandleFunc(shared.StatusPath, handlers.StatusHandler)
	mux.HandleFunc(shared.StatusPath[:len(shared.StatusPath)-1], handlers.StatusHandler)

	// Registrations
	mux.HandleFunc(shared.RegistrationsPath, handlers.RegistrationsHandler)
	mux.HandleFunc(
		shared.RegistrationsPath[:len(shared.RegistrationsPath)-1],
		handlers.RegistrationsHandler,
	)

	// Registrations with ID
	mux.HandleFunc(shared.RegistrationsPath+"{id}", handlers.RegistrationsHandlerWithID)

	// Default, serves the web page
	mux.HandleFunc("/", handlers.DefaultHandler)

	// mux.HandleFunc("/dashboard/v1/registrations/", listRegistrationsHandler)
	// mux.HandleFunc("/dashboard/v1/registrations/{id}", registrationsHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
