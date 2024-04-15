package handlers

import (
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/handlers/dashboards"
	"assignment-2/internal/http/handlers/notifications"
	"assignment-2/internal/http/handlers/registrations"
	"assignment-2/internal/http/handlers/status"
	"encoding/json"
	"log"
	"net/http"
)

type siteMap struct {
	Help      string             `json:"help"`
	Endpoints []inhouse.Endpoint `json:"siteMap"`
}

// SiteMap
// Site map for the server. Contains all the endpoints and their descriptions.
var SiteMap = siteMap{
	Help: "This is the default handler for the server. " +
		"Maybe you typed the wrong path or are looking for the web page. Go to '/' for the web page.",
	Endpoints: []inhouse.Endpoint{},
}

// Init
// Initializes the site map with all the endpoints from the different handlers.
func Init() {
	endpointsFromRegistrations := registrations.GetEndpointStructs()
	endpointsFromDashboards := dashboards.GetEndpointStructs()
	endpointsFromNotifications := notifications.GetEndpointStructs()
	endpointsFromStatus := status.GetEndpointStructs()

	SiteMap.Endpoints = append(SiteMap.Endpoints, endpointsFromRegistrations...)
	SiteMap.Endpoints = append(SiteMap.Endpoints, endpointsFromDashboards...)
	SiteMap.Endpoints = append(SiteMap.Endpoints, endpointsFromNotifications...)
	SiteMap.Endpoints = append(SiteMap.Endpoints, endpointsFromStatus...)
}

// DefaultHandler
// Default handler for the server. Redirects to the web page.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// If the request is for the root path, redirect to the web page
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/web", http.StatusSeeOther)
		return
	}

	// Else, return the site map
	w.Header().Set("content-type", "application/json")
	marshaledSiteMap, err := json.MarshalIndent(SiteMap, "", "\t")
	if err != nil {
		log.Println("Error during JSON encoding: " + err.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshaledSiteMap)
	if err != nil {
		log.Println("Failed to write response: " + err.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
