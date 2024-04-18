package notifications

import (
	"assignment-2/internal/datasources/firebase"
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

func FindNotifications(event string) ([]requests.Notification, error) {
	notifications, err := firebase.GetAllDocuments[requests.Notification](firebase.NotificationCollection)
	if err != nil {
		log.Println("Error while trying to receive notification from db: ", err.Error())
		return nil, err
	}

	for _, notification := range notifications {
		if notification.Event == event {
			notifications = append(notifications, notification)
		}
	}
	return notifications, nil
}

/*
FindNotificationsByCountry returns all notifications for a specific event and country.
*/
func FindNotificationsByCountry(event string, country string) ([]requests.Notification, error) {
	var foundNotifications []requests.Notification
	notifications, err := firebase.GetAllDocuments[requests.Notification](firebase.NotificationCollection)
	if err != nil {
		log.Println("Error while trying to receive notification from db: ", err.Error())
		return nil, err
	}

	for _, notification := range notifications {
		fmt.Printf("notification country: %v\n", notification.Country)
		fmt.Printf("Given country: %v\n", country)
		if notification.Event == event && (notification.Country == country || notification.Country == "") {
			fmt.Printf("inside if\n")
			foundNotifications = append(foundNotifications, notification)
		}
	}
	return foundNotifications, nil
}

// InvokeNotification sends a request to the URL of the notification with the content of the notification as the body.
func InvokeNotification(notification requests.Notification) {
	// Update the notification with the current time
	currentTime := time.Now()
	notification.LastInvoke = &currentTime

	err := firebase.UpdateDocument[requests.Notification](
		notification, notification.ID,
		firebase.NotificationCollection,
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
	r, err3 := http.NewRequest(http.MethodGet, notification.Url, reader)
	if err3 != nil {
		log.Println("Error while creating request: " + err3.Error())
		return
	}
	_, err4 := utils.Client.Do(r)
	if err4 != nil {
		log.Println("Error while sending request: " + err4.Error())
		return
	}
}
