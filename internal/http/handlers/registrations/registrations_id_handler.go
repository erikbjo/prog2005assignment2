package registrations

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/db"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/http/handlers/notifications"
	"assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Implemented methods for the endpoint with ID
var implementedMethodsWithID = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
}

// Endpoint for managing registrations with a specific ID
var registrationsEndpointWithID = inhouse.Endpoint{
	Path:        constants.RegistrationsPath + "{id}",
	Methods:     implementedMethodsWithID,
	Description: "This endpoint is used to manage registrations with a specific ID.",
}

// HandlerWithID handles the /registrations/v1/registrations/{id} path.
func HandlerWithID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleRegistrationsGetRequestWithID(w, r)
	case http.MethodPut:
		handleRegistrationsPutRequestWithID(w, r)
	case http.MethodDelete:
		handleRegistrationsDeleteRequestWithID(w, r)

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

func handleRegistrationsGetRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the registration with the provided ID
	dashboard, err2 := db.GetDocument[requests.DashboardConfig](
		id,
		db.DashboardCollection,
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
		dashboard,
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

func handleRegistrationsPutRequestWithID(w http.ResponseWriter, r *http.Request) {
	var update requests.DashboardConfig

	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&update); err != nil {
		log.Println(constants.ErrJsonDecode + err.Error())
		http.Error(w, constants.ErrJsonDecode, http.StatusBadRequest)
		return
	}

	update.ID = id
	update.LastChange = time.Now()

	err3 := db.UpdateDocument[requests.DashboardConfig](
		update, id,
		db.DashboardCollection,
	)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}

	// Check if any notifications are registered for the event
	foundNotifications, err4 := notifications.FindNotificationsByCountry(
		requests.EventChange,
		update.IsoCode,
	)
	if err4 != nil {
		log.Println(constants.ErrNotificationsGetDocFromDB, err4.Error())
		http.Error(w, constants.ErrNotificationsGetDocFromDB, http.StatusInternalServerError)
		return
	}

	// If found, invoke the notifications
	if len(foundNotifications) > 0 {
		for _, n := range foundNotifications {
			notifications.InvokeNotification(n)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleRegistrationsDeleteRequestWithID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the registration with the provided ID
	dashboard, err3 := db.GetDocument[requests.DashboardConfig](
		id,
		db.DashboardCollection,
	)
	if err3 != nil {
		switch err3.Error() {
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
		log.Println(constants.ErrDBGetDoc + err3.Error())
		return
	}

	err2 := db.DeleteDocument(id, db.DashboardCollection)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any notifications are registered for the event
	foundNotifications, err4 := notifications.FindNotificationsByCountry(
		requests.EventDelete,
		dashboard.IsoCode,
	)
	if err4 != nil {
		log.Println(constants.ErrNotificationsGetDocFromDB, err4.Error())
		http.Error(w, constants.ErrNotificationsGetDocFromDB, http.StatusInternalServerError)
		return
	}

	// If found, invoke the notifications
	if len(foundNotifications) > 0 {
		for _, n := range foundNotifications {
			notifications.InvokeNotification(n)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
