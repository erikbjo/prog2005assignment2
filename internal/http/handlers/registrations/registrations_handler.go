package registrations

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/db"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/http/handlers/notifications"
	"assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Implemented methods for the endpoint without ID
var implementedMethodsWithoutID = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
}

// Endpoint for managing registrations without a specific ID
var registrationsEndpointWithoutID = inhouse.Endpoint{
	Path:        constants.RegistrationsPath,
	Methods:     implementedMethodsWithoutID,
	Description: "This endpoint is used to manage registrations.",
}

/*
HandlerWithoutID handles the /dashboard/v1/registrations path.
*/
func HandlerWithoutID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleRegistrationsGetRequest(w, r)
	case http.MethodHead:
		// Advanced Task: Implement the HEAD method functionality (only return the header, not the body).
		handleRegistrationsHeadRequest(w, r)
	case http.MethodPost:
		handleRegistrationsPostRequest(w, r)

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

/*
handleRegistrationsGetRequest handles the GET request for the /dashboard/v1/registrations path.
*/
func handleRegistrationsGetRequest(w http.ResponseWriter, r *http.Request) {

	// Get the all dashboard config documents
	allDocuments, err2 := db.GetAllDocuments[requests.DashboardConfig](db.DashboardCollection)
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

/*
handleRegistrationsHeadRequest handles the HEAD request for the /dashboard/v1/registrations path.
*/
func handleRegistrationsHeadRequest(w http.ResponseWriter, r *http.Request) {

	// Get all dashboard config documents to get content length
	allDocuments, err2 := db.GetAllDocuments[requests.DashboardConfig](db.DashboardCollection)
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
		allDocuments,
		"",
		"\t",
	)
	if err3 != nil {
		log.Println("Error during JSON encoding: " + err3.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Set response headers
	headers := map[string]string{
		"Date":           time.Now().Format(time.RFC1123),
		"Content-Type":   r.Header.Get("Content-Type"),
		"Connection":     r.Header.Get("Connection"),
		"Content-Length": strconv.Itoa(len(marshaled)),
	}

	fmt.Println(headers)
	// Set response headers
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(http.StatusOK)
}

/*
handleRegistrationsPostRequest handles the POST request for the /dashboard/v1/registrations path.
*/
func handleRegistrationsPostRequest(w http.ResponseWriter, r *http.Request) {

	var content requests.DashboardConfig

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&content)
	if err != nil {
		log.Println("Error while decoding json: ", err.Error())
		http.Error(w, "Error while decoding json.", http.StatusBadRequest)
		return
	}

	content.LastChange = time.Now()
	content.ID = utils.GenerateRandomID()

	// Save the DashboardConfig to the database
	err2 := db.AddDocument[requests.DashboardConfig](content, db.DashboardCollection)
	if err2 != nil {
		http.Error(w, "Error while trying to add document.", http.StatusInternalServerError)
	}

	// Check if any notifications are registered for the event
	foundNotifications, err3 := notifications.FindNotificationsByCountry(
		requests.EventRegister,
		content.IsoCode,
	)
	if err3 != nil {
		log.Println("Error while trying to find notifications: ", err3.Error())
		http.Error(w, "Error while trying to find notifications.", http.StatusInternalServerError)
		return
	}

	// If found, invoke the notifications
	if len(foundNotifications) > 0 {
		for _, n := range foundNotifications {
			notifications.InvokeNotification(n)
		}
	}

	marshaled, err4 := json.MarshalIndent(
		registrationResponse{ID: content.ID, LastChange: content.LastChange},
		"",
		"\t",
	)
	if err4 != nil {
		log.Println("Error during JSON encoding: " + err4.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Set the status code to 201 Created
	w.WriteHeader(http.StatusCreated)
	// Write the JSON to the response
	_, err5 := w.Write(marshaled)
	if err5 != nil {
		log.Println("Failed to write response: " + err5.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
