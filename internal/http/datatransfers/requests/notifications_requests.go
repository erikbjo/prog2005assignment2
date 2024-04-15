package requests

import "time"

const (
	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventDelete   = "DELETE"
	EventInvoke   = "INVOKE"
)

type Notification struct {
	ID      string    `json:"id"`
	Url     string    `json:"url"`
	Country string    `json:"country"`
	Event   string    `json:"event"`
	Time    time.Time `json:"time"`
}
