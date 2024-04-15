package notifications

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
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
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received request to get notification with ID %s\n", id)

	// Get the notification with the provided ID
	notification, err2 := firebase.GetDocument[requests.Notification](
		id,
		firebase.NotificationCollection,
	)
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
		notification,
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

func handleNotificationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Received request to delete notification with ID ", id)

	err2 := firebase.DeleteDocument(id, firebase.NotificationCollection)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}
