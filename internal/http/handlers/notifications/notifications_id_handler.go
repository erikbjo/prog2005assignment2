package notifications

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/http/datatransfers/inhouse"
	"fmt"
	"net/http"
)

// Implemented methods for the endpoint with ID
var implementedMethodsWithID = []string{
	http.MethodGet,
	http.MethodDelete,
}

// Endpoint for managing notifications with a specific ID
var notificationsEndpointWithID = inhouse.Endpoint{
	Path:        constants.NotificationsPath + "{id}",
	Methods:     implementedMethodsWithID,
	Description: "This endpoint is used to manage notifications with a specific ID.",
}

// HandlerWithID handles the /notifications/{id} path.
func HandlerWithID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleNotificationsGetRequestWithID(w, r)
	case http.MethodDelete:
		handleNotificationsDeleteRequestWithID(w, r)

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

func handleNotificationsGetRequestWithID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "GET request not implemented", http.StatusNotImplemented)
}

func handleNotificationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "DELETE request not implemented", http.StatusNotImplemented)
}
