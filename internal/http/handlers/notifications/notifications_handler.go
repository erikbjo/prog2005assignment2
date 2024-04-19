package notifications

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/db"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Implemented methods for the endpoint
var implementedMethodsWithoutID = []string{
	http.MethodGet,
	http.MethodPost,
}

// Endpoint for managing webhooks for event notifications.
var notificationsEndpointWithoutID = inhouse.Endpoint{
	Path:        constants.NotificationsPath,
	Methods:     implementedMethodsWithoutID,
	Description: "Endpoint for managing webhooks for event notifications.",
}

// HandlerWithoutID handles the /notifications path.
// It currently supports GET, POST and DELETE requests.
// Endpoint for managing webhooks for event notifications.
func HandlerWithoutID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
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
	// Get all notification documents from db
	allDocuments, err2 := firebase.GetAllDocuments[requests.Notification](firebase.NotificationCollection)
	if err2 != nil {
		http.Error(
			w,
			"Error while trying to receive document from db.",
			http.StatusInternalServerError,
		)
		log.Println("Error while trying to receive document from db: ", err2.Error())
		return
	}
	if len(allDocuments) > 0 {
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
	} else {
		http.Error(w, "No documents found", http.StatusNoContent)
	}
}

func handleNotificationsPostRequest(w http.ResponseWriter, r *http.Request) {
	var content requests.Notification

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&content); err != nil {
		log.Println("Error while decoding json: ", err.Error())
	}

	// Checks if event in body is isValid
	if isValidEvent(content.Event) == false {
		http.Error(w, "Invalid event type provided", http.StatusBadRequest)
		return
	}

	content.ID = utils.GenerateRandomID()

	// Save the Notification to the database
	err2 := db.AddDocument[requests.Notification](content, db.NotificationCollection)
	if err2 != nil {
		http.Error(w, "Error while trying to add document.", http.StatusInternalServerError)
	}

	// Return the ID of the saved Notification
	// Marshal the status object to JSON
	marshaled, err3 := json.MarshalIndent(
		notificationResponse{Id: content.ID},
		"",
		"\t",
	)
	if err3 != nil {
		log.Println("Error during JSON encoding: " + err3.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Set the status code to 201 Created
	w.WriteHeader(http.StatusCreated)
	// Write the JSON to the response
	_, err4 := w.Write(marshaled)
	if err4 != nil {
		log.Println("Failed to write response: " + err4.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

/*
isValidEvent checks if the event is a isValid event type.
*/
func isValidEvent(event string) bool {
	for _, validEvent := range requests.ImplementedEvents {
		if event == validEvent {
			return true
		}
	}
	return false
}
