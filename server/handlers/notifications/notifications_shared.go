package notifications

import "assignment-2/server/shared"

// GetEndpointStructs returns the endpoints for the registrations handler. One with an ID and one without.
func GetEndpointStructs() []shared.Endpoint {
	return []shared.Endpoint{notificationsEndpointWithoutID, notificationsEndpointWithID}
}
