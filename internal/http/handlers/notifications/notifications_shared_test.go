package notifications

import (
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/datatransfers/requests"
	"testing"
	"time"
)

func TestFindNotifications(t *testing.T) {
	type args struct {
		event string
	}
	tests := []struct {
		name                string
		args                args
		onlyWantedEventType bool
		wantError           bool
	}{
		{
			name: "FoundTestFindNotifications",
			args: args{
				event: "REGISTER",
			},
			onlyWantedEventType: true,
			wantError:           false,
		},
		{
			name: "InvalidTestFindNotifications",
			args: args{
				event: "INVALID_EVENT",
			},
			onlyWantedEventType: false,
			wantError:           true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := FindNotifications(tt.args.event)

				// Check if the number of notifications returned is as expected
				if !(len(got) == 0) && func() bool {
					for _, n := range got {
						if n.Event != tt.args.event {
							return false
						}
					}
					return true
				}() != tt.onlyWantedEventType {
					t.Errorf(
						"FindNotifications() found unwanted event type in list = %v, wanted type %v",
						got,
						tt.args.event,
					)
				}
				// Check if the error is as expected
				if (err != nil) != tt.wantError {
					t.Errorf("FindNotifications() error = %v, wantErr %v", err, tt.wantError)
					return
				}
			},
		)
	}
}

func TestFindNotificationsByCountry(t *testing.T) {
	type args struct {
		event   string
		country string
	}
	tests := []struct {
		name                          string
		args                          args
		onlyWantedEventTypeAndCountry bool
		wantError                     bool
	}{
		{
			name: "FoundTestFindNotificationsByCountry",
			args: args{
				event:   "REGISTER",
				country: "SE",
			},
			onlyWantedEventTypeAndCountry: true,
			wantError:                     false,
		},
		{
			name: "InvalidEventTestFindNotificationsByCountry",
			args: args{
				event:   "INVALID_EVENT",
				country: "SE",
			},
			onlyWantedEventTypeAndCountry: false,
			wantError:                     true,
		},
		{
			name: "InvalidCountryTestFindNotificationsByCountry",
			args: args{
				event:   "REGISTER",
				country: "",
			},
			onlyWantedEventTypeAndCountry: true,
			wantError:                     false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := FindNotificationsByCountry(tt.args.event, tt.args.country)

				// Check if the number of notifications returned is as expected
				if !(len(got) == 0) && func() bool {
					for _, n := range got {
						// Check if event matches and country matches (including empty string)
						if n.Event == tt.args.event && (n.Country == tt.args.country || n.Country == "") {
							return true
						}
					}
					// No match found in the loop
					return false
				}() != tt.onlyWantedEventTypeAndCountry {
					t.Errorf(
						"FindNotificationsByCountry() found unwanted event type or country in list = %v, "+
							"wanted type %v and country %v", got,
						tt.args.event, tt.args.country,
					)
				}
				// Check if the error is as expected
				if (err != nil) != tt.wantError {
					t.Errorf(
						"FindNotificationsByCountry() error = %v, wantErr %v",
						err,
						tt.wantError,
					)
					return
				}
			},
		)
	}
}

func TestInvokeNotification(t *testing.T) {
	type args struct {
		notification requests.Notification
	}
	tests := []struct {
		name            string
		args            args
		wantTimeUpdated bool
	}{
		{
			name: "ValidTestInvokeNotification",
			args: args{
				notification: requests.Notification{
					ID:         "testID",
					Url:        "http://localhost:8080",
					Event:      "REGISTER",
					Country:    "SE",
					LastInvoke: &time.Time{},
				},
			},
			wantTimeUpdated: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				timeBefore := time.Now()
				*tt.args.notification.LastInvoke = timeBefore

				_ = firebase.AddDocument[requests.Notification](
					tt.args.notification,
					firebase.NotificationCollection,
				)
				InvokeNotification(tt.args.notification)
				timeAfter := tt.args.notification.LastInvoke
				if (timeAfter != &timeBefore) != tt.wantTimeUpdated {
					t.Errorf(
						"InvokeNotification() did not update time, time before update: %v, time after update: %v",
						timeBefore, timeAfter,
					)
				}
			},
		)
	}
}
