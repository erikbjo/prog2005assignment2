// Package handlers
// This file contains shared variables and constants used by the handlers
package handlers

import (
	"assignment-2/server/shared"
	"net/http"
	"time"
)

// Current API URLs, changes if testing
var (
	currentRestCountriesApi = shared.RestCountriesApi
	currentCurrencyApi      = shared.CurrencyApi
)

// startTime is the time the server started
var startTime = time.Now()

// client is the HTTP client used to send requests
var client = &http.Client{
	Timeout: 3 * time.Second,
}

// SetStubsForTesting Use self-hosted stubs for testing
func SetStubsForTesting() {
	// TODO: Implement the stubs to test
	currentRestCountriesApi = shared.TestRestCountriesApi
	currentCurrencyApi = shared.TestCurrencyApi
}
