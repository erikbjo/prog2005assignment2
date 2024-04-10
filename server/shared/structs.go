package shared

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

type DashboardConfig struct {
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

type SiteMap struct {
	Help      string     `json:"help"`
	Endpoints []Endpoint `json:"siteMap"`
}

type Endpoint struct {
	Path        string   `json:"path"`
	Methods     []string `json:"methods"`
	Description string   `json:"description"`
}
