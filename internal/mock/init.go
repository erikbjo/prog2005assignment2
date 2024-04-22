package mock

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/db"
	"assignment-2/internal/mock/stubs"
	"assignment-2/internal/utils"
	"log"
	"net/http"
	"time"
)

func InitForTesting() {
	setStubsForTesting()
	// Initialize Firebase for testing
	db.InitializeForTesting()
	createTestHttpServer()
}

var localhost = "http://localhost:" + utils.GetTestPort()
var server http.Server

// setStubsForTesting Use self-hosted stubs for testing
func setStubsForTesting() {
	utils.CurrentRestCountriesApi = localhost + constants.TestRestCountriesApi
	utils.CurrentCurrencyApi = localhost + constants.TestCurrencyApi
	utils.CurrentMeteoApi = localhost + constants.TestMeteoApi
}

func createTestHttpServer() {
	if server.Addr == "" {
		port := utils.GetTestPort()
		server := http.Server{Addr: ":" + port}

		http.HandleFunc(constants.TestRestCountriesApi, stubs.RestCountriesHandler)
		http.HandleFunc(constants.TestCurrencyApi, stubs.CurrencyHandler)
		http.HandleFunc(constants.TestMeteoApi, stubs.MeteoHandler)

		go func() {
			// If 8001 is in use, sleep for 1 second and try again
			// This is to avoid the error "listen tcp :8001: bind: address already in use"
			// when running multiple tests
			// Hideous
			err := server.ListenAndServe()
			if err != nil {
				if err.Error() == "listen tcp :8001: bind: address already in use" {
					log.Println("Port 8001 is in use, sleeping for 5 seconds and trying again")
					time.Sleep(5 * time.Second)
					err = server.ListenAndServe()
					if err != nil {
						log.Fatalf("Failed to start http server: %v", err)
					}
				} else {
					log.Fatalf("Failed to start http server: %v", err)
				}
			}
		}()
	}
}

func TeardownAfterTesting() {
	// Currently not used as Go creates a new instance of the server for each test
}
