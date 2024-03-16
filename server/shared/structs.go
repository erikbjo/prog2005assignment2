package shared

// A Status struct to hold the status of the server, including the status of the APIs and the version
// of the server.
type Status struct {
	CountriesAPI   int     `json:"countries_api"`
	MeteoAPI       int     `json:"meteo_api"`
	CurrencyAPI    int     `json:"currency_api"`
	NotificationDB int     `json:"notification_db"`
	Webhooks       int     `json:"webhooks"`
	Version        string  `json:"version"`
	Uptime         float64 `json:"uptime"`
}
