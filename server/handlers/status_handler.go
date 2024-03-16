package handlers

import (
	"assignment-2/server/shared"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

// Current API URLs, changes if testing
var (
	currentRestCountriesApi = shared.RestCountriesApi
	currentCurrencyApi      = shared.CurrencyApi
)

// StartTime is the time the server started
var startTime = time.Now()

// client is the HTTP client used to send requests
var client = &http.Client{
	Timeout: 3 * time.Second,
}

// StatusHandler
// Status handler for the server. Returns the status of the server and the APIs it relies on.
// Currently only supports GET requests.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{http.MethodGet}

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r, false)

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

// handleStatusGetRequest handles the GET request for the /status path.
// It returns the status of the server and the APIs it relies on.
func handleStatusGetRequest(w http.ResponseWriter, r *http.Request, testing bool) {
	if testing {
		currentRestCountriesApi = shared.TestRestCountriesApi
		currentCurrencyApi = shared.TestCurrencyApi
	}

	// Create a new status object
	// TODO: Implement the MeteoAPI, NotificationDB, Webhooks
	currentStatus := shared.Status{
		CountriesAPI:   getStatusCode(currentRestCountriesApi, w),
		MeteoAPI:       http.StatusNotImplemented,
		CurrencyAPI:    getStatusCode(currentCurrencyApi, w),
		NotificationDB: http.StatusNotImplemented,
		Webhooks:       http.StatusNotImplemented,
		Version:        shared.Version,
		Uptime:         math.Round(time.Since(startTime).Seconds()),
	}

	// Marshal the status object to JSON
	marshaledStatus, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		log.Println("Error during JSON encoding: " + err.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response
	_, err = w.Write(marshaledStatus)
	if err != nil {
		log.Println("Failed to write response: " + err.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// getStatusCode returns the status code of the given URL.
// If the URL is not reachable, it returns 503.
func getStatusCode(url string, w http.ResponseWriter) int {
	// Send a GET request to the URL
	resp, err := client.Get(url)
	if err != nil {
		// If there is an error, return 503
		return http.StatusServiceUnavailable
	}

	// Return the status code
	return resp.StatusCode
}
