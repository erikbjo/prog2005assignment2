package mock

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/mock/stubs"
	"assignment-2/internal/utils"
	"log"
	"net/http"
	"time"
)

func InitForTesting() {
	setStubsForTesting()
	// Initialize Firebase for testing
	firebase.InitializeForTesting()
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
					log.Println("Port 8001 is in use, sleeping for 2 seconds and trying again")
					time.Sleep(2 * time.Second)
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
	// cmd := exec.Command(
	// 	"curl",
	// 	"-v",
	// 	"-X",
	// 	"DELETE",
	// 	"http://localhost:8888/emulator/v1/projects/prog2005-assignment-2-c2e5c/databases/(" +
	//	"default)/documents",
	// )
	// if err := cmd.Run(); err != nil {
	// 	log.Fatalf("Failed to stop Firestore emulator: %v", err)
	// }

	// Close http server
	// err := server.Close()
	// if err != nil {
	// 	log.Fatalf("Failed to shutdown http server: %v", err)
	// }

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
