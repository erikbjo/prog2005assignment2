package notifications

import (
	"assignment-2/server/shared"
	"fmt"
	"net/http"
)

var implementedMethods = []string{
	http.MethodGet,
	http.MethodPost,
}

var notificationsEndpoint = shared.Endpoint{
	Path:        shared.NotificationsPath,
	Methods:     implementedMethods,
	Description: "Endpoint for managing webhooks for event notifications.",
}

func GetEndpointStructs() []shared.Endpoint {
	return []shared.Endpoint{notificationsEndpoint}
}

// Handler handles the /notifications path.
// It currently supports GET, POST and DELETE requests.
// Endpoint for managing webhooks for event notifications.
func Handler(w http.ResponseWriter, r *http.Request) {
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
				implementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}

}

// NotificationsHandlerWithID handles the /notifications/{id} path.
func NotificationsHandlerWithID(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{
		http.MethodGet,
		http.MethodDelete,
	}

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
				implementedMethods,
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

func handleNotificationsGetRequestWithID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "GET request not implemented", http.StatusNotImplemented)
}

func handleNotificationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "DELETE request not implemented", http.StatusNotImplemented)
}
