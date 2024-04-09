package handlers

import (
	"assignment-2/server/handlers/registrations"
	"assignment-2/server/shared"
	"encoding/json"
	"log"
	"net/http"
)

var SiteMap = shared.SiteMap{
	Help:      "This is the default handler for the server. It redirects to the web page.",
	Endpoints: []shared.Endpoint{},
}

func Init() {
	var1 := registrations.GetEndpointStructs()
	SiteMap.Endpoints = append(SiteMap.Endpoints, var1...)
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
