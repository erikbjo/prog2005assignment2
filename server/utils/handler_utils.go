// Package utils
// This file contains shared variables and constants used by the handlers
package utils

import (
	"assignment-2/server/shared"
	"net/http"
	"time"
)

// Current API URLs, changes if testing
var (
	CurrentRestCountriesApi = shared.RestCountriesApi
	CurrentCurrencyApi      = shared.CurrencyApi
)

// StartTime is the time the server started
var StartTime = time.Now()

// Client is the HTTP Client used to send requests
var Client = &http.Client{
	Timeout: 3 * time.Second,
}

// SetStubsForTesting Use self-hosted stubs for testing
func SetStubsForTesting() {
	// TODO: Implement the stubs to test
	CurrentRestCountriesApi = shared.TestRestCountriesApi
	CurrentCurrencyApi = shared.TestCurrencyApi
}
