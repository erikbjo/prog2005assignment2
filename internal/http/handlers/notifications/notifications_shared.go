package notifications

import (
	"assignment-2/internal/db"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type notificationResponse struct {
	Id string `json:"id"`
}

// GetEndpointStructs returns the endpoints for the registrations handler. One with an ID and one without.
func GetEndpointStructs() []inhouse.Endpoint {
	return []inhouse.Endpoint{notificationsEndpointWithoutID, notificationsEndpointWithID}
}

/*
FindNotifications returns all notifications for a specific event without any other conditions.
*/
func FindNotifications(event string) ([]requests.Notification, error) {
	var foundNotifications []requests.Notification

	if !isValidEvent(event) {
		return nil, fmt.Errorf("invalid event type: %v", event)
	}

	notifications, err := db.GetAllDocuments[requests.Notification](db.NotificationCollection)
	if err != nil {
		log.Println("Error while trying to receive notification from db: ", err.Error())
		return nil, err
	}

	for _, notification := range notifications {
		if notification.Event == event {
			foundNotifications = append(foundNotifications, notification)
		}
	}
	return foundNotifications, nil
}

/*
FindNotificationsByCountry returns all notifications for a specific event and country as condition.
*/
func FindNotificationsByCountry(event string, country string) ([]requests.Notification, error) {
	var foundNotifications []requests.Notification

	if !isValidEvent(event) {
		return nil, fmt.Errorf("invalid event type: %v", event)
	}

	notifications, err := db.GetAllDocuments[requests.Notification](db.NotificationCollection)
	if err != nil {
		log.Println("Error while trying to receive notification from db: ", err.Error())
		return nil, err
	}

	for _, notification := range notifications {
		if notification.Event == event && (notification.Country == country || notification.Country == "") {
			foundNotifications = append(foundNotifications, notification)
		}
	}
	return foundNotifications, nil
}

// InvokeNotification invokes the notification by sending a request to the URL of the notification with the content of
// the notification as the body.
func InvokeNotification(notification requests.Notification) {
	// Update the notification with the current time
	currentTime := time.Now()
	notification.LastInvoke = &currentTime

	err := db.UpdateDocument[requests.Notification](
		notification, notification.ID,
		db.NotificationCollection,
	)
	if err != nil {
		log.Println("Error while updating document: " + err.Error())
		return
	}

	// Marshal the status object to JSON
	marshaled, err2 := json.MarshalIndent(
		notification,
		"",
		"\t",
	)
	if err2 != nil {
		log.Println("Error during JSON marshaling: " + err2.Error())
		return
	}

	reader := bytes.NewReader(marshaled)
	r, err3 := http.NewRequest(http.MethodPost, notification.Url, reader)
	if err3 != nil {
		log.Println("Error while creating request: " + err3.Error())
		return
	}

	// Sets header
	r.Header.Set("Content-Type", "application/json")

	_, err4 := utils.Client.Do(r)
	if err4 != nil {
		log.Println("Error while sending request: " + err4.Error())
		return
	}
}
