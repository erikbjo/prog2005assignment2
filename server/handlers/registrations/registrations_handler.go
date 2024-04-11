package registrations

import (
	"assignment-2/db"
	"assignment-2/server/shared"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Implemented methods for the endpoint without ID
var implementedMethodsWithoutID = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
}

// Endpoint for managing registrations without a specific ID
var registrationsEndpointWithoutID = shared.Endpoint{
	Path:        shared.RegistrationsPath,
	Methods:     implementedMethodsWithoutID,
	Description: "This endpoint is used to manage registrations.",
}

// HandlerWithoutID handles the /registrations/v1/registrations path.
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

	// Get the all dashboard config documents
	allDocuments, err2 := db.GetAllDocuments(w, r, db.DashboardCollection)
	if err2 != nil {
		http.Error(w, "Error while trying to receive document from db.", http.StatusInternalServerError)
		log.Println("Error while trying to receive document from db: ", err2.Error())
		return
	}

	// Marshal the status object to JSON
	marshaled, err3 := json.MarshalIndent(
		allDocuments,
		"",
		"\t",
	)
	if err3 != nil {
		log.Println("Error during JSON encoding: " + err3.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response
	_, err4 := w.Write(marshaled)
	if err4 != nil {
		log.Println("Failed to write response: " + err4.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
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

	/*
		// Read and parse the body
		validRequest, err := checkValidityOfResponseBody(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	*/

	// log.Println(validRequest)

	// Save the DashboardConfig to the database
	id, configMap, err2 := db.AddDashboardConfigDocument(w, r, db.DashboardCollection)
	if err2 != nil {
		http.Error(w, "Error while trying to add document.", http.StatusInternalServerError)
	}

	// Return the ID of the saved DashboardConfig
	// TODO: Implement returning the ID of the saved DashboardConfig
	// Marshal the status object to JSON
	marshaled, err3 := json.MarshalIndent(
		shared.RegistrationResponse{ID: id, LastChange: configMap.LastChange},
		"",
		"\t",
	)
	if err3 != nil {
		log.Println("Error during JSON encoding: " + err3.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response
	_, err4 := w.Write(marshaled)
	if err4 != nil {
		log.Println("Failed to write response: " + err4.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
