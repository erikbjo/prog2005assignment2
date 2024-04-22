package notifications

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/db"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Implemented methods for the endpoint with ID
var implementedMethodsWithID = []string{
	http.MethodGet,
	http.MethodDelete,
}

// Endpoint for managing notifications with a specific ID
var notificationsEndpointWithID = inhouse.Endpoint{
	Path:        constants.NotificationsPath + "{id}",
	Methods:     implementedMethodsWithID,
	Description: "This endpoint is used to manage notifications with a specific ID.",
}

// HandlerWithID handles the /notifications/{id} path.
func HandlerWithID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleNotificationsGetRequestWithID(w, r)
	case http.MethodDelete:
		handleNotificationsDeleteRequestWithID(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethodsWithID,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleNotificationsGetRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the notification with the provided ID
	notification, err2 := db.GetDocument[requests.Notification](
		id,
		db.NotificationCollection,
	)
	if err2 != nil {
		switch err2.Error() {
		case constants.ErrIDInvalid:
			http.Error(w, constants.ErrIDInvalid, http.StatusBadRequest)
		case constants.ErrDBDocNotFound:
			http.Error(w, constants.ErrDBDocNotFound, http.StatusNoContent)
		default:
			http.Error(
				w,
				constants.ErrDBGetDoc,
				http.StatusInternalServerError,
			)
		}

		log.Println(constants.ErrDBGetDoc + err2.Error())
		return
	}

	// Marshal the status object to JSON
	marshaled, err3 := json.MarshalIndent(
		notification,
		"",
		"\t",
	)
	if err3 != nil {
		log.Println(constants.ErrJsonMarshal + err3.Error())
		http.Error(w, constants.ErrJsonMarshal, http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response
	_, err4 := w.Write(marshaled)
	if err4 != nil {
		log.Println(constants.ErrWriteResponse + err4.Error())
		http.Error(w, constants.ErrWriteResponse, http.StatusInternalServerError)
		return
	}
}

func handleNotificationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := db.DeleteDocument(id, db.NotificationCollection)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
