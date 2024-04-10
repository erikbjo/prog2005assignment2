package dashboards

import (
	"assignment-2/db"
	"assignment-2/server/shared"
	"fmt"
	"log"
	"net/http"
)

// Implemented methods for the endpoint
var implementedMethods = []string{
	http.MethodGet,
}

// Endpoint for managing dashboards
var dashboardsEndpoint = shared.Endpoint{
	Path:        shared.DashboardsPath,
	Methods:     implementedMethods,
	Description: "Endpoint for managing dashboards.",
}

// GetEndpointStructs returns the endpoint struct for the dashboards endpoint.
func GetEndpointStructs() []shared.Endpoint {
	return []shared.Endpoint{dashboardsEndpoint}
}

// HandlerWithID handles the /dashboard/v1/dashboards path.
// It currently only supports GET requests
func HandlerWithID(w http.ResponseWriter, r *http.Request) {
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleDashboardsGetRequest(w, r)

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

// handleDashboardsGetRequest handles the GET request for the /dashboard/v1/dashboards path.
// It is used to retrieve the populated dashboards.
func handleDashboardsGetRequest(w http.ResponseWriter, r *http.Request) {
	if len(r.PathValue("id")) == 0 {
		http.Error(w, "No document ID was provided.", http.StatusBadRequest)
	} else {
		mp, err := db.DisplayDocument(w, r, db.DashboardCollection)
		if err != nil {
			log.Println("Error while trying to display dashboard document: ", err.Error())
			http.Error(
				w,
				"Error while trying to display dashboard document.",
				http.StatusInternalServerError,
			)
		}
		log.Println("Received request with map: ", mp)
	}
}
