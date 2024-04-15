package responses

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
