package mock

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/utils"
)

func InitForTesting() {
	setStubsForTesting()
	// Initialize Firebase for testing
	firebase.InitializeForTesting()
}

// setStubsForTesting Use self-hosted stubs for testing
func setStubsForTesting() {
	// TODO: Implement the stubs to mock
	utils.CurrentRestCountriesApi = constants.TestRestCountriesApi
	utils.CurrentCurrencyApi = constants.TestCurrencyApi
	utils.CurrentMeteoApi = constants.TestMeteoApi
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
