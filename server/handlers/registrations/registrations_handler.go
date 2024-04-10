package registrations

import (
	"assignment-2/server/shared"
	"fmt"
	"log"
	"net/http"
)

var implementedMethodsWithoutID = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
}

var registrationsEndpointWithoutID = shared.Endpoint{
	Path:        shared.RegistrationsPath,
	Methods:     implementedMethodsWithoutID,
	Description: "This endpoint is used to manage registrations.",
}

func HandlerWithoutID(w http.ResponseWriter, r *http.Request) {
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
				implementedMethodsWithoutID,
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
	validRequest, err := checkValidityOfResponseBody(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(validRequest)

	// Save the DashboardConfig to the database
	// TODO: Implement saving the DashboardConfig to the database

	// Return the ID of the saved DashboardConfig
	// TODO: Implement returning the ID of the saved DashboardConfig

	http.Error(w, "PUT request not implemented", http.StatusNotImplemented)
}
