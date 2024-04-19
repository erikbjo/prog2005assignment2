package mock

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/mock/stubs"
	"assignment-2/internal/utils"
	"log"
	"net/http"
	"os/exec"
)

func InitForTesting() {
	setStubsForTesting()
	// Initialize Firebase for testing
	firebase.InitializeForTesting()
	createTestHttpServer()
}

var localhost = "http://localhost:" + utils.GetPort()

// setStubsForTesting Use self-hosted stubs for testing
func setStubsForTesting() {
	utils.CurrentRestCountriesApi = localhost + constants.TestRestCountriesApi
	utils.CurrentCurrencyApi = localhost + constants.TestCurrencyApi
	utils.CurrentMeteoApi = localhost + constants.TestMeteoApi
}

func createTestHttpServer() {
	port := utils.GetPort()

	http.HandleFunc(constants.TestRestCountriesApi, stubs.RestCountriesHandler)
	http.HandleFunc(constants.TestCurrencyApi, stubs.CurrencyHandler)
	http.HandleFunc(constants.TestMeteoApi, stubs.MeteoHandler)

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			panic(err)
		}
	}()

}

func TeardownAfterTesting() {
	cmd := exec.Command(
		"curl",
		"-v",
		"-X",
		"DELETE",
		"http://localhost:8888/emulator/v1/projects/prog2005-assignment-2-c2e5c/databases/(default)/documents",
	)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to stop Firestore emulator: %v", err)
	}

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
