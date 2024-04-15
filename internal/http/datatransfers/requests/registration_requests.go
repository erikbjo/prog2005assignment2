package requests

import "time"

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
