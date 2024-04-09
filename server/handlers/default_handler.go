package handlers

import (
	"assignment-2/server/shared"
	"net/http"
)

var siteMap = shared.SiteMap{
	Help:      "This is the default handler for the server. It redirects to the web page.",
	Endpoints: []shared.Endpoint{},
}

// DefaultHandler
// Default handler for the server. Redirects to the web page.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// If the request is for the root path, redirect to the web page
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/web", http.StatusSeeOther)
		return
	}

}
