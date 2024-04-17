package server

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/handlers"
	"assignment-2/internal/http/handlers/dashboards"
	"assignment-2/internal/http/handlers/notifications"
	"assignment-2/internal/http/handlers/registrations"
	"assignment-2/internal/http/handlers/status"
	"assignment-2/internal/utils"
	"log"
	"net/http"
)

// Start
/*
Start the server on the port specified in the environment variable PORT. If PORT is not set, the default port 8080 is used.
*/
func Start() {
	// Initialization of firebase database
	firebase.Initialize()

	// Firebase client closes at the end of this function
	defer firebase.Close()

	// Get the port from the environment variable, or use the default port
	port := utils.GetPort()

	// Using mux to handle /'s and parameters
	mux := http.NewServeMux()

	// Initialize the site map
	handlers.Init()

	// Set up handler endpoints, with and without trailing slash
	// Status
	mux.HandleFunc(constants.StatusPath, status.Handler)
	mux.HandleFunc(constants.StatusPath[:len(constants.StatusPath)-1], status.Handler)

	// Registrations
	mux.HandleFunc(constants.RegistrationsPath, registrations.HandlerWithoutID)
	mux.HandleFunc(
		constants.RegistrationsPath[:len(constants.RegistrationsPath)-1],
		registrations.HandlerWithoutID,
	)

	// Registrations with ID
	mux.HandleFunc(constants.RegistrationsPath+"{id}", registrations.HandlerWithID)

	// Dashboards
	mux.HandleFunc(constants.DashboardsPath+"{id}", dashboards.HandlerWithID)

	// Notifications
	mux.HandleFunc(constants.NotificationsPath, notifications.HandlerWithoutID)
	mux.HandleFunc(
		constants.NotificationsPath[:len(constants.NotificationsPath)-1],
		notifications.HandlerWithoutID,
	)

	// Notifications with ID
	mux.HandleFunc(constants.NotificationsPath+"{id}", notifications.HandlerWithID)

	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/web/", http.StripPrefix("/web/", fs))

	// Default, redirect to /web/
	mux.HandleFunc("/", handlers.DefaultHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
