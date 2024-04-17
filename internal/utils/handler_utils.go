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

const FirestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"

// StartTime is the time the server started
var StartTime = time.Now()

// Client is the HTTP Client used to send requests
var Client = &http.Client{
	Timeout: 3 * time.Second,
}

// setStubsForTesting Use self-hosted stubs for testing
func setStubsForTesting() {
	// TODO: Implement the stubs to mock
	CurrentRestCountriesApi = constants.TestRestCountriesApi
	CurrentCurrencyApi = constants.TestCurrencyApi
	CurrentMeteoApi = constants.TestMeteoApi
}

func SetupForTesting() {

	// firebase.Initialize()
	// Set stubs for testing
	setStubsForTesting()
}

func TeardownAfterTesting() {
	// // Clean up any resources after all tests have been executed
	// firebase.Close()
	//
	// // firebase emulators:stop
	// cmd := exec.Command("firebase", "emulators:stop")
	// if err := cmd.Run(); err != nil {
	// 	log.Fatalf("Failed to stop Firestore emulator: %v", err)
	// }
	// os.Exit(0)
}
