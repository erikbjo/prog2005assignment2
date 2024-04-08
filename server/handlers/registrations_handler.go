package handlers

import (
	"assignment-2/server/shared"
	"assignment-2/server/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
	}

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleRegistrationsGetRequest(w, r)
	case http.MethodHead:
		// Advanced Task: Implement the HEAD method functionality (only return the header, not the body).
		handleRegistrationsHeadRequest(w, r)
	case http.MethodPost:
		handleRegistrationsPostRequest(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}

}

func RegistrationsHandlerWithID(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
	}

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleRegistrationsGetRequestWithID(w, r)
	case http.MethodPut:
		handleRegistrationsPutRequestWithID(w, r)
	case http.MethodDelete:
		handleRegistrationsDeleteRequestWithID(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}

}

func handleRegistrationsGetRequest(w http.ResponseWriter, r *http.Request) {
	// Pseudocode
	// Get all registrations from the database
	// Return the registrations
	// If there is an error, return an error message

	http.Error(w, "GET request not implemented", http.StatusNotImplemented)
}

// Advanced Task: Implement the HEAD method functionality (only return the header, not the body).
func handleRegistrationsHeadRequest(w http.ResponseWriter, r *http.Request) {
	// Pseudocode
	// Get all registrations from the database
	// Return the headers only
	// If there is an error, return an error message

	http.Error(w, "HEAD request not implemented", http.StatusNotImplemented)

}

func handleRegistrationsPostRequest(w http.ResponseWriter, r *http.Request) {
	// Pseudocode
	// Parse the body
	// Decode the body into a DashboardConfig struct
	// Save the DashboardConfig to the database
	// Return the ID of the saved DashboardConfig
	// If there is an error, return an error message

	// Read and parse the body
	body, err := getDashboardConfigFromBody(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Received request to create registration with country ", body.Country)

	// Save the DashboardConfig to the database
	// TODO: Implement saving the DashboardConfig to the database

	// Return the ID of the saved DashboardConfig
	// TODO: Implement returning the ID of the saved DashboardConfig

	http.Error(w, "PUT request not implemented", http.StatusNotImplemented)
}

func handleRegistrationsGetRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Received request to get registration with ID ", id)

	// Get the registration with the provided ID
	// TODO: Implement getting the registration with the provided ID

	http.Error(w, "GET request not implemented", http.StatusNotImplemented)
}

func handleRegistrationsPutRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := getDashboardConfigFromBody(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Received request to update registration with country ", body.Country)

	log.Println("Received request to update registration with ID ", id)

	http.Error(w, "POST request not implemented", http.StatusNotImplemented)
}

func handleRegistrationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Received request to delete registration with ID ", id)

	http.Error(w, "DELETE request not implemented", http.StatusNotImplemented)
}

func getDashboardConfigFromBody(w http.ResponseWriter, r *http.Request) (
	shared.DashboardConfig,
	error,
) {
	var dashboardConfig shared.DashboardConfig

	// Read and parse the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body: ", err)
		// Note: We don't return the error here because we want to return a generic error message
		return dashboardConfig, fmt.Errorf("error reading request body")
	}

	if len(body) == 0 {
		log.Println("Empty request body")
		return dashboardConfig, fmt.Errorf("empty request body")
	}

	// Decode the body into a DashboardConfig struct
	err = json.Unmarshal(body, &dashboardConfig)
	if err != nil {
		log.Println("Error decoding request body: ", err)
		return dashboardConfig, fmt.Errorf("error decoding request body")
	}

	return dashboardConfig, nil
}
