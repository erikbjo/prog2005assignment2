package registrations

import (
	`assignment-2/internal/http/datatransfers/inhouse`
	`assignment-2/internal/http/datatransfers/requests`
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type registrationResponse struct {
	ID         string    `json:"id"`
	LastChange time.Time `json:"lastChange"`
}

// GetEndpointStructs returns the endpoints for the registrations handler. One with an ID and one without.
func GetEndpointStructs() []inhouse.Endpoint {
	return []inhouse.Endpoint{registrationsEndpointWithoutID, registrationsEndpointWithID}
}

func checkValidityOfResponseBody(w http.ResponseWriter, r *http.Request) (
	bool,
	error,
) {
	var dashboardConfig requests.DashboardConfig
	var copyOfBody = r.Body

	// Read and parse the body
	body, err := io.ReadAll(copyOfBody)
	if err != nil {
		log.Println("Error reading request body: ", err)
		// Note: We don't return the error here because we want to return a generic error message
		return false, fmt.Errorf("error reading request body")
	}

	if len(body) == 0 {
		log.Println("Empty request body")
		return false, fmt.Errorf("empty request body")
	}

	// Decode the body into a DashboardConfig struct
	err = json.Unmarshal(body, &dashboardConfig)
	if err != nil {
		log.Println("Error decoding request body: ", err)
		return false, fmt.Errorf("error decoding request body")
	}

	// Additional checks can be added here to validate the body

	return true, nil
}