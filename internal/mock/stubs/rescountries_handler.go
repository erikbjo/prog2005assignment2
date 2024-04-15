package stubs

import (
	"fmt"
	"log"
	"net/http"
)

func RestCountriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("content-type", "application/json")
		output := ParseFile("./mock/resources/restcountries_no.json")

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
