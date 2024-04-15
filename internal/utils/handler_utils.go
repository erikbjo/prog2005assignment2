// Package utils
// This file contains shared variables and constants used by the handlers
package utils

import (
	"assignment-2/internal/constants"
	"net/http"
	"time"
)

// Current API URLs, changes if testing
var (
	CurrentRestCountriesApi = constants.RestCountriesApi
	CurrentCurrencyApi      = constants.CurrencyApi
	CurrentMeteoApi         = constants.MeteoApi
)

// StartTime is the time the server started
var StartTime = time.Now()

// Client is the HTTP Client used to send requests
var Client = &http.Client{
	Timeout: 3 * time.Second,
}

// SetStubsForTesting Use self-hosted stubs for testing
func SetStubsForTesting() {
	// TODO: Implement the stubs to mock
	CurrentRestCountriesApi = constants.TestRestCountriesApi
	CurrentCurrencyApi = constants.TestCurrencyApi
	CurrentMeteoApi = constants.TestMeteoApi
}
