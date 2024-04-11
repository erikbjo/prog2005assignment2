package shared

import "time"

// A Status struct to hold the status of the server, including the status of the APIs and the version
// of the server.
type Status struct {
	CountriesAPI   int     `json:"countries_api"`
	MeteoAPI       int     `json:"meteo_api"`
	CurrencyAPI    int     `json:"currency_api"`
	DashboardDB    int     `json:"dashboard_db"`
	NotificationDB int     `json:"notification_db"`
	Webhooks       int     `json:"webhooks"`
	Version        string  `json:"version"`
	Uptime         float64 `json:"uptime"`
}

type RegistrationResponse struct {
	ID         string    `json:"id"`
	LastChange time.Time `json:"lastChange"`
}

type DashboardConfig struct {
	ID       string `json:"id"`
	Country  string `json:"country"`
	IsoCode  string `json:"isoCode"`
	Features struct {
		Temperature      bool     `json:"temperature"`
		Precipitation    bool     `json:"precipitation"`
		Capital          bool     `json:"capital"`
		Coordinates      bool     `json:"coordinates"`
		Population       bool     `json:"population"`
		Area             bool     `json:"area"`
		TargetCurrencies []string `json:"targetCurrencies"`
	} `json:"features"`
	LastChange time.Time `json:"lastChange"`
}

type Dashboard struct {
	Country  string `json:"country"`
	IsoCode  string `json:"isoCode"`
	Features struct {
		Temperature   float64 `json:"temperature"`
		Precipitation float64 `json:"precipitation"`
		Capital       string  `json:"capital"`
		Coordinates   struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"coordinates"`
		Population       int     `json:"population"`
		Area             float64 `json:"area"`
		TargetCurrencies struct {
			EUR float64 `json:"EUR"`
			USD float64 `json:"USD"`
			SEK float64 `json:"SEK"`
		} `json:"targetCurrencies"`
	} `json:"features"`
	LastRetrieval string `json:"lastRetrieval"`
}

type MeteoForecastResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Elevation            float64 `json:"elevation"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Hourly               struct {
		Time          []interface{} `json:"time"`
		Temperature2M []interface{} `json:"temperature_2m"`
	} `json:"hourly,omitempty"`
	HourlyUnits struct {
		Temperature2M string `json:"temperature_2m"`
	} `json:"hourly_units,omitempty"`
}

type SiteMap struct {
	Help      string     `json:"help"`
	Endpoints []Endpoint `json:"siteMap"`
}

type Endpoint struct {
	Path        string   `json:"path"`
	Methods     []string `json:"methods"`
	Description string   `json:"description"`
}
