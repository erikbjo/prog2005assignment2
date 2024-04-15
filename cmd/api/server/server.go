package server

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/handlers"
	"assignment-2/internal/http/handlers/dashboards"
	notifications2 "assignment-2/internal/http/handlers/notifications"
	registrations2 "assignment-2/internal/http/handlers/registrations"
	"assignment-2/internal/http/handlers/status"
	stubs2 "assignment-2/internal/mock/stubs"
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
	mux.HandleFunc(constants.RegistrationsPath, registrations2.HandlerWithoutID)
	mux.HandleFunc(
		constants.RegistrationsPath[:len(constants.RegistrationsPath)-1],
		registrations2.HandlerWithoutID,
	)

	// Registrations with ID
	mux.HandleFunc(constants.RegistrationsPath+"{id}", registrations2.HandlerWithID)

	// Dashboards
	mux.HandleFunc(constants.DashboardsPath+"{id}", dashboards.HandlerWithID)

	// Notifications
	mux.HandleFunc(constants.NotificationsPath, notifications2.HandlerWithoutID)
	mux.HandleFunc(
		constants.NotificationsPath[:len(constants.NotificationsPath)-1],
		notifications2.HandlerWithoutID,
	)

	// Notifications with ID
	mux.HandleFunc(constants.NotificationsPath+"{id}", notifications2.HandlerWithID)

	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/web/", http.StripPrefix("/web/", fs))

	// Default, redirect to /web/
	mux.HandleFunc("/", handlers.DefaultHandler)

	// Set up stubs for testing
	mux.HandleFunc(constants.TestMeteoApi, stubs2.MeteoHandler)
	mux.HandleFunc(constants.TestRestCountriesApi, stubs2.RestCountriesHandler)
	mux.HandleFunc(constants.TestCurrencyApi, stubs2.CurrencyHandler)

	// mux.HandleFunc("/dashboard/v1/registrations/", listRegistrationsHandler)
	// mux.HandleFunc("/dashboard/v1/registrations/{id}", registrationsHandler)

	// mux.HandleFunc("/dbTest/", db.HandleDB)
	// mux.HandleFunc("/dbTest/{id}/", db.HandleDB)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
