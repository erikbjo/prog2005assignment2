package stubs

import (
	"fmt"
	"log"
	"net/http"
)

func RestCountriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// cwd, err := os.Getwd()
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	return
		// }

		// Print the current working directory
		// fmt.Println("Current working directory:", cwd)

		w.Header().Add("content-type", "application/json")
		// Hideous way to get the path to the mock resources, but it works for now
		output := ParseFile("../../../mock/resources/restcountries_no.json")

		_, err := fmt.Fprint(w, string(output))
		if err != nil {
			log.Println("Error while trying to display the Restcountries API stub: ", err.Error())
			http.Error(
				w,
				"Error while trying to display the Restcountries API stub.",
				http.StatusInternalServerError,
			)
			return
		}
		break
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}
