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
	ID         string         `json:"id"`
	Country    string         `json:"country"`
	IsoCode    string         `json:"isoCode"`
	Features   ConfigFeatures `json:"features"`
	LastChange time.Time      `json:"lastChange"`
}

type ConfigFeatures struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

type Dashboard struct {
	Country       string            `json:"country"`
	IsoCode       string            `json:"isoCode"`
	Features      DashboardFeatures `json:"features"`
	LastRetrieval string            `json:"lastRetrieval"`
}

type DashboardFeatures struct {
	Temperature      float64               `json:"temperature"`
	Precipitation    float64               `json:"precipitation"`
	Capital          []string              `json:"capital"`
	Coordinates      Coordinates           `json:"coordinates"`
	Population       int                   `json:"population"`
	Area             float64               `json:"area"`
	TargetCurrencies map[string]float64    `json:"targetCurrencies"`
	Currencies       map[string]Currencies `json:"currencies"`
}

type Currencies struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ResponseFromRestcountries struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Cca2       string                `json:"cca2"`
	Currencies map[string]Currencies `json:"currencies"`
	Capital    []string              `json:"capital"`
	Latlng     []float64             `json:"latlng"`
	Area       float64               `json:"area"`
	Population int                   `json:"population"`
}

type ResponseFromMeteo struct {
	CurrentUnits struct {
		Time          string `json:"time"`
		Interval      string `json:"interval"`
		Temperature2M string `json:"temperature_2m"`
		Precipitation string `json:"precipitation"`
	} `json:"current_units"`
	Current struct {
		Time          string  `json:"time"`
		Interval      int     `json:"interval"`
		Temperature2M float64 `json:"temperature_2m"`
		Precipitation float64 `json:"precipitation"`
	} `json:"current"`
}

type RestCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Cca2       string                `json:"cca2"`
	Currencies map[string]Currencies `json:"currencies"`
	Capital    []string              `json:"capital"`
	Latlng     []int                 `json:"latlng"`
	Area       int                   `json:"area"`
	Population int                   `json:"population"`
}

type ResponseFromCurrency struct {
	Result   string             `json:"result"`
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"rates"`
}

type Notification struct {
	Id      string    `json:"id"`
	Country string    `json:"country"`
	Event   string    `json:"event"`
	Time    time.Time `json:"time"`
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
