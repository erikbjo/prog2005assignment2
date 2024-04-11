package stubs

import (
	"fmt"
	"log"
	"net/http"
)

func MeteoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("content-type", "application/json")
		output := ParseFile("./test/resources/meteo_norway.json")

		_, err := fmt.Fprint(w, string(output))
		if err != nil {
			log.Println("Error while trying to display the Meteo API stub: ", err.Error())
			http.Error(
				w,
				"Error while trying to display the Meteo API stub.",
				http.StatusInternalServerError,
			)
			return
		}
		break
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}
