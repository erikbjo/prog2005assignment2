package handlers

import (
	"assignment-2/db"
	"fmt"
	"log"
	"net/http"
)

// DashboardsHandlerWithID handles the /dashboard/v1/dashboards path.
// It currently only supports GET requests
func DashboardsHandlerWithID(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{
		http.MethodGet,
	}

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
		err := db.DisplayDocument(w, r, db.DashboardCollection)
		if err != nil {
			log.Println("Error while trying to display dashboard document: ", err.Error())
			http.Error(w, "Error while trying to display dashboard document.", http.StatusInternalServerError)
		}
	}
}
