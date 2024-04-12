package registrations

import (
	"assignment-2/db"
	"assignment-2/server/shared"
	"assignment-2/server/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Implemented methods for the endpoint with ID
var implementedMethodsWithID = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
}

// Endpoint for managing registrations with a specific ID
var registrationsEndpointWithID = shared.Endpoint{
	Path:        shared.RegistrationsPath + "{id}",
	Methods:     implementedMethodsWithID,
	Description: "This endpoint is used to manage registrations with a specific ID.",
}

// HandlerWithID handles the /registrations/v1/registrations/{id} path.
func HandlerWithID(w http.ResponseWriter, r *http.Request) {
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
				implementedMethodsWithID,
			), http.StatusNotImplemented,
		)
		return
	}

}

func handleRegistrationsGetRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received request to get registration with ID %s\n", id)

	// Get the registration with the provided ID
	// TODO: Implement getting the registration with the provided ID
	dashboard, err2 := db.GetDashboardConfigDocument(id, db.DashboardCollection)
	if err2 != nil {
		http.Error(
			w,
			"Error while trying to receive document from db.",
			http.StatusInternalServerError,
		)
		log.Println("Error while trying to receive document from db: ", err2.Error())
		return
	}

	// Marshal the status object to JSON
	marshaled, err3 := json.MarshalIndent(
		dashboard,
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

func handleRegistrationsPutRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: This still leads to EOF later in function
	/*
		body, err := checkValidityOfResponseBody(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	log.Println("Received request to update registration with ID ", id)

	err2 := db.UpdateDashboardConfigDocument(w, r, db.DashboardCollection)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}
}

func handleRegistrationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Received request to delete registration with ID ", id)

	err2 := db.DeleteDocument(w, r, db.DashboardCollection)
	if err2 != nil {
		http.Error(w, "Error while trying to delete document.", http.StatusInternalServerError)
		return
	}
}
