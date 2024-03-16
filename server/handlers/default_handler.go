package handlers

import (
	"net/http"
)

// DefaultHandler
// Default handler for the server. Redirects to the web page.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file from the web directory
	http.ServeFile(w, r, "web/index.html")
}
