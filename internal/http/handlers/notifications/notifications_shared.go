package notifications

import (
	"assignment-2/internal/http/datatransfers/inhouse"
	"time"
)

type Notification struct {
	Id      string    `json:"id"`
	Country string    `json:"country"`
	Event   string    `json:"event"`
	Time    time.Time `json:"time"`
}

// GetEndpointStructs returns the endpoints for the registrations handler. One with an ID and one without.
func GetEndpointStructs() []inhouse.Endpoint {
	return []inhouse.Endpoint{notificationsEndpointWithoutID, notificationsEndpointWithID}
}
