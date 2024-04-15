package requests

const (
	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventDelete   = "DELETE"
	EventInvoke   = "INVOKE"
)

type Notification struct {
	Id      string `json:"id"`
	Url     string `json:"url"`
	Country string `json:"country"`
	Event   string `json:"event"`
}
