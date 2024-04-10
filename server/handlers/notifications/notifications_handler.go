package notifications

import (
	"assignment-2/server/shared"
	"fmt"
	"net/http"
)

// Implemented methods for the endpoint
var implementedMethodsWithoutID = []string{
	http.MethodGet,
	http.MethodPost,
}

// Endpoint for managing webhooks for event notifications.
var notificationsEndpointWithoutID = shared.Endpoint{
	Path:        shared.NotificationsPath,
	Methods:     implementedMethodsWithoutID,
	Description: "Endpoint for managing webhooks for event notifications.",
}

// HandlerWithoutID handles the /notifications path.
// It currently supports GET, POST and DELETE requests.
// Endpoint for managing webhooks for event notifications.
func HandlerWithoutID(w http.ResponseWriter, r *http.Request) {
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleNotificationsGetRequest(w, r)
	case http.MethodPost:
		handleNotificationsPostRequest(w, r)

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

func handleNotificationsGetRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "GET request not implemented", http.StatusNotImplemented)
}

func handleNotificationsPostRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "POST request not implemented", http.StatusNotImplemented)
}
