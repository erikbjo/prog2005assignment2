package registrations

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

// HandlerWithoutID handles the /dashboard/v1/registrations path.
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

func handleRegistrationsGetRequest(w http.ResponseWriter, r *http.Request) {
	// Pseudocode
	// Get all registrations from the database
	// Return the registrations
	// If there is an error, return an error message

	// Get the all dashboard config documents
	allDocuments, err2 := firebase.GetAllDocuments(firebase.DashboardCollection)
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

	// Write the JSON to the response
	_, err4 := w.Write(marshaled)
	if err4 != nil {
		log.Println("Failed to write response: " + err4.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// Returns only the headers for the request to retrieve all dashboards configs
func handleRegistrationsHeadRequest(w http.ResponseWriter, r *http.Request) {

	// Get all dashboard config documents to get content length
	allDocuments, err2 := firebase.GetAllDocuments(firebase.DashboardCollection)
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

func handleRegistrationsPostRequest(w http.ResponseWriter, r *http.Request) {
	// Pseudocode
	// Parse the body
	// Decode the body into a DashboardConfig struct
	// Save the DashboardConfig to the database
	// Return the ID of the saved DashboardConfig
	// If there is an error, return an error message

	/*
		// Read and parse the body
		validRequest, err := checkValidityOfResponseBody(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	*/

	// log.Println(validRequest)

	var content requests.DashboardConfig

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&content); err != nil {
		log.Println("Error while decoding json: ", err.Error())
	}

	content.LastChange = time.Now()
	content.ID = utils.GenerateRandomID()

	// Save the DashboardConfig to the database
	err2 := firebase.AddDocument[requests.DashboardConfig](content, firebase.DashboardCollection)
	if err2 != nil {
		http.Error(w, "Error while trying to add document.", http.StatusInternalServerError)
	}

	// Return the ID of the saved DashboardConfig
	// TODO: Implement returning the ID of the saved DashboardConfig
	// Marshal the status object to JSON
	marshaled, err3 := json.MarshalIndent(
		registrationResponse{ID: content.ID, LastChange: content.LastChange},
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
