package registrations

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/http/handlers/notifications"
	"assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Implemented methods for the endpoint with ID
var implementedMethodsWithID = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
}

// Endpoint for managing registrations with a specific ID
var registrationsEndpointWithID = inhouse.Endpoint{
	Path:        constants.RegistrationsPath + "{id}",
	Methods:     implementedMethodsWithID,
	Description: "This endpoint is used to manage registrations with a specific ID.",
}

// HandlerWithID handles the /registrations/v1/registrations/{id} path.
func HandlerWithID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
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
	dashboard, err2 := firebase.GetDocument[requests.DashboardConfig](
		id,
		firebase.DashboardCollection,
	)
	if err2 != nil {
		switch err2.Error() {
		case "no valid ID was provided":
			http.Error(w, "No valid ID was provided", http.StatusBadRequest)
		case "document not found in collection":
			http.Error(w, "Document not found in collection", http.StatusNoContent)
		default:
			http.Error(
				w, "Error while trying to receive document from db.", http.StatusInternalServerError,
			)
		}
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
	var update requests.DashboardConfig

	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&update); err != nil {
		log.Println("Error while decoding json: ", err.Error())
		http.Error(w, "Error while decoding json.", http.StatusBadRequest)
		return
	}

	update.ID = id
	update.LastChange = time.Now()

	log.Println("Received request to update registration with ID ", id)

	err3 := firebase.UpdateDocument[requests.DashboardConfig](
		update, id,
		firebase.DashboardCollection,
	)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}

	// Check if any notifications are registered for the event
	foundNotifications, err4 := notifications.FindNotificationsByCountry(requests.EventChange, update.IsoCode)
	if err4 != nil {
		log.Println("Error while trying to find notifications: ", err4.Error())
		http.Error(w, "Error while trying to find notifications.", http.StatusInternalServerError)
		return
	}

	// If found, invoke the notifications
	if len(foundNotifications) > 0 {
		for _, n := range foundNotifications {
			notifications.InvokeNotification(n)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleRegistrationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the registration with the provided ID
	dashboard, err3 := firebase.GetDocument[requests.DashboardConfig](
		id,
		firebase.DashboardCollection,
	)
	if err3 != nil {
		switch err3.Error() {
		case "no valid ID was provided":
			http.Error(w, "No valid ID was provided", http.StatusBadRequest)
		case "document not found in collection":
			http.Error(w, "Document not found in collection", http.StatusNoContent)
		default:
			http.Error(
				w, "Error while trying to receive document from db.", http.StatusInternalServerError,
			)
		}
		log.Println("Error while trying to receive document from db: ", err3.Error())
		return
	}

	log.Println("Received request to delete registration with ID ", id)

	err2 := firebase.DeleteDocument(id, firebase.DashboardCollection)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any notifications are registered for the event
	foundNotifications, err4 := notifications.FindNotificationsByCountry(requests.EventDelete, dashboard.IsoCode)
	if err4 != nil {
		log.Println("Error while trying to find notifications: ", err4.Error())
		http.Error(w, "Error while trying to find notifications.", http.StatusInternalServerError)
		return
	}

	// If found, invoke the notifications
	if len(foundNotifications) > 0 {
		for _, n := range foundNotifications {
			notifications.InvokeNotification(n)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
