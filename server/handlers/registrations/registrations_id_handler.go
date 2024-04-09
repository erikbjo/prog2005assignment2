package registrations

import (
	"assignment-2/server/shared"
	"assignment-2/server/utils"
	"fmt"
	"log"
	"net/http"
)

var implementedMethodsWithID = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
}

var registrationsEndpointWithID = shared.Endpoint{
	Path:        "/registrations/{id}",
	Methods:     implementedMethodsWithID,
	Description: "This endpoint is used to manage registrations with a specific ID.",
}

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

	body, err := checkValidityOfResponseBody(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(body)

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
