package handlers

import (
	"fmt"
	"net/http"
)

func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{
		http.MethodGet,
		http.MethodPost,
	}

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleRegistrationsGetRequest(w, r)
	case http.MethodPut:
		handleRegistrationsPutRequest(w, r)

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
	case http.MethodPost:
		handleRegistrationsPostRequestWithID(w, r)
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
	http.Error(w, "GET request not implemented", http.StatusNotImplemented)
}

func handleRegistrationsPutRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "PUT request not implemented", http.StatusNotImplemented)
}

func handleRegistrationsGetRequestWithID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "GET request not implemented", http.StatusNotImplemented)
}

func handleRegistrationsPostRequestWithID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "POST request not implemented", http.StatusNotImplemented)
}

func handleRegistrationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "DELETE request not implemented", http.StatusNotImplemented)
}
