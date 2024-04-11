package dashboards

import (
	"assignment-2/db"
	"assignment-2/server/shared"
	"assignment-2/server/utils"
	"encoding/json"
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
	id, err := utils.GetIDFromRequest(r)

	mp, err := db.GetDashboardConfigDocument(id, db.DashboardCollection)
	if err != nil {
		log.Println("Error while trying to display dashboard document: ", err.Error())
		http.Error(
			w,
			"Error while trying to display dashboard document.",
			http.StatusInternalServerError,
		)
	}

	// mp is map[string]interface {} type
	// convert it to shared.DashboardConfig type

	dashboard := mp

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dashboard)
	if err != nil {
		log.Println("Error while trying to encode dashboard document: ", err.Error())
		http.Error(
			w,
			"Error while trying to encode dashboard document.",
			http.StatusInternalServerError,
		)
	}
}
