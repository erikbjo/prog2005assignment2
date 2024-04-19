package handlers

import (
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/handlers/dashboards"
	"assignment-2/internal/http/handlers/notifications"
	"assignment-2/internal/http/handlers/registrations"
	"assignment-2/internal/http/handlers/status"
	"encoding/json"
	"github.com/russross/blackfriday"
	"log"
	"net/http"
	"os"
)

type siteMap struct {
	Help      string             `json:"help"`
	Endpoints []inhouse.Endpoint `json:"siteMap"`
}

// SiteMap
// Site map for the server. Contains all the endpoints and their descriptions.
var SiteMap = siteMap{
	Help: "This is the default handler for the server. " +
		"Maybe you typed the wrong path or are looking for the web page. " +
		"Go to '/' to read the README.md file." +
		"ID is required for endpoints with path ending in {id}.",
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
		// Read the contents of the README.md file
		readme, err := os.ReadFile("README.md")
		if err != nil {
			http.Error(w, "Failed to read README.md", http.StatusInternalServerError)
			return
		}

		html := blackfriday.MarkdownCommon(readme)

		// Set the Content-Type header to indicate that this is HTML
		w.Header().Set("Content-Type", "text/html")

		// Custom styles for the HTML to make it look better
		customStyles := `
            <style>
                pre {background-color: #f4f4f4;}
				body {padding: 10px;}
            </style>
        `

		// Append the custom styles to the HTML
		htmlWithStyles := append([]byte(customStyles), html...)

		// Write the HTML response
		_, err = w.Write(htmlWithStyles)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	} else {
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
}
