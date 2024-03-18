package handlers

import (
	"fmt"
	"net/http"
)

// NotificationsHandler handles the /notifications path.
// It currently supports GET, POST and DELETE requests.
// Endpoint for managing webhooks for event notifications.
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
	}

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleNotificationsGetRequest(w, r)
	case http.MethodPost:
		handleNotificationsPostRequest(w, r)
	case http.MethodDelete:
		handleNotificationsDeleteRequest(w, r)

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

func handleNotificationsDeleteRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "DELETE request not implemented", http.StatusNotImplemented)
}