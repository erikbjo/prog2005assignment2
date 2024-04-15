package notifications

import (
	"assignment-2/internal/http/datatransfers/inhouse"
)

type notificationResponse struct {
	Id string `json:"id"`
}

// GetEndpointStructs returns the endpoints for the registrations handler. One with an ID and one without.
func GetEndpointStructs() []inhouse.Endpoint {
	return []inhouse.Endpoint{notificationsEndpointWithoutID, notificationsEndpointWithID}
}
