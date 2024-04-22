package registrations

import (
	"assignment-2/internal/http/datatransfers/inhouse"
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
