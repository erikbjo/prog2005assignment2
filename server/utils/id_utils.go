package utils

import (
	"fmt"
	"log"
	"net/http"
)

// GetIDFromRequest Get the ID from the request
func GetIDFromRequest(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if id == "" {
		log.Println("ID not provided")
		return "", fmt.Errorf("ID not provided")
	}

	return id, nil
}
