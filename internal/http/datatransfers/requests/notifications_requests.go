package requests

import "time"

const (
	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventDelete   = "DELETE"
	EventInvoke   = "INVOKE"
)

// Slice of currently implemented event types for notifications
var ImplementedEvents = []string{EventRegister, EventChange, EventDelete, EventInvoke}

type Notification struct {
	ID      string    `json:"id"`
	Url     string    `json:"url"`
	Country string    `json:"country"`
	Event   string    `json:"event"`
	Time    time.Time `json:"time"`
}
