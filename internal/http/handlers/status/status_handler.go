package status

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

// A Status struct to hold the status of the server, including the status of the APIs and the version
// of the server.
type status struct {
	CountriesAPI   int    `json:"countries_api"`
	MeteoAPI       int    `json:"meteo_api"`
	CurrencyAPI    int    `json:"currency_api"`
	DashboardDB    int    `json:"dashboard_db"`
	NotificationDB int    `json:"notification_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         int    `json:"uptime"`
}

// implementedMethods is a list of the implemented HTTP methods for the status endpoint.
var implementedMethods = []string{http.MethodGet}

// statusEndpoint is the endpoint for checking the status of the server and the APIs it relies on.
var statusEndpoint = inhouse.Endpoint{
	Path:        constants.StatusPath,
	Methods:     implementedMethods,
	Description: "Endpoint for checking the status of the server and the APIs it relies on.",
}

// GetEndpointStructs returns the endpoint for the status handler.
func GetEndpointStructs() []inhouse.Endpoint {
	return []inhouse.Endpoint{statusEndpoint}
}

// Handler
// Status handler for the server. Returns the status of the server and the APIs it relies on.
// Currently only supports GET requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r)

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
func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {
	// Create a new status object
	// TODO: Implement the Webhooks
	// TODO: Implement firebase testing/mocking
	currentStatus := status{
		CountriesAPI:   getStatusCode(utils.CurrentRestCountriesApi, w),
		MeteoAPI:       getStatusCode(utils.CurrentMeteoApi, w),
		CurrencyAPI:    getStatusCode(utils.CurrentCurrencyApi, w),
		DashboardDB:    firebase.GetStatusCodeOfCollection(w, firebase.DashboardCollection),
		NotificationDB: firebase.GetStatusCodeOfCollection(w, firebase.NotificationCollection),
		Webhooks:       http.StatusNotImplemented,
		Version:        constants.Version,
		Uptime:         int(math.Round(time.Since(utils.StartTime).Seconds())),
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
	switch url {
	case utils.CurrentRestCountriesApi:
		url = url + "all"
	case utils.CurrentCurrencyApi:
		url = url + "nok"
	case utils.CurrentMeteoApi:
		url = url + "?latitude=60.7957&longitude=10.6915"
	}

	// Send a GET request to the URL
	resp, err := utils.Client.Get(url)
	if err != nil {
		// If there is an error, return 503
		return http.StatusServiceUnavailable
	}

	// Return the status code
	return resp.StatusCode
}
