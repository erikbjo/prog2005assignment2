package server

import (
	"assignment-2/db"
	"assignment-2/server/handlers"
	"assignment-2/server/handlers/dashboards"
	"assignment-2/server/handlers/notifications"
	"assignment-2/server/handlers/registrations"
	"assignment-2/server/handlers/status"
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

	// Firebase client closes at the end of this function
	defer db.Close()

	// Get the port from the environment variable, or use the default port
	port := utils.GetPort()

	// Using mux to handle /'s and parameters
	mux := http.NewServeMux()

	// Initialize the site map
	handlers.Init()

	// Set up handler endpoints, with and without trailing slash
	// Status
	mux.HandleFunc(shared.StatusPath, status.StatusHandler)
	mux.HandleFunc(shared.StatusPath[:len(shared.StatusPath)-1], status.StatusHandler)

	// Registrations
	mux.HandleFunc(shared.RegistrationsPath, registrations.HandlerWithoutID)
	mux.HandleFunc(
		shared.RegistrationsPath[:len(shared.RegistrationsPath)-1],
		registrations.HandlerWithoutID,
	)

	// Registrations with ID
	mux.HandleFunc(shared.RegistrationsPath+"{id}", registrations.HandlerWithID)

	// Dashboards
	mux.HandleFunc(shared.DashboardsPath, dashboards.DashboardsHandlerWithID)
	mux.HandleFunc(
		shared.DashboardsPath[:len(shared.DashboardsPath)-1],
		dashboards.DashboardsHandlerWithID,
	)

	// Notifications
	mux.HandleFunc(shared.NotificationsPath, notifications.NotificationsHandler)
	mux.HandleFunc(
		shared.NotificationsPath[:len(shared.NotificationsPath)-1],
		notifications.NotificationsHandler,
	)

	// Notifications with ID
	mux.HandleFunc(shared.NotificationsPath+"{id}", notifications.NotificationsHandlerWithID)

	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/web/", http.StripPrefix("/web/", fs))

	// Default, redirect to /web/
	mux.HandleFunc("/", handlers.DefaultHandler)

	// mux.HandleFunc("/dashboard/v1/registrations/", listRegistrationsHandler)
	// mux.HandleFunc("/dashboard/v1/registrations/{id}", registrationsHandler)

	// mux.HandleFunc("/dbTest/", db.HandleDB)
	// mux.HandleFunc("/dbTest/{id}/", db.HandleDB)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
