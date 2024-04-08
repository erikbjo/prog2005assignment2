package utils

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetIDFromRequest(r *http.Request) (int, error) {
	id := r.PathValue("id")
	if id == "" {
		log.Println("ID not provided")
		return -1, fmt.Errorf("ID not provided")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Invalid ID, must be an integer")
		return -1, fmt.Errorf("invalid ID, must be an integer")
	}

	return idInt, nil
}
