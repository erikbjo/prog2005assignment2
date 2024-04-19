package notifications

import (
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/mock"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testNotification = requests.Notification{
	Url:     "http://localhost:8080/test/",
	Country: "NO",
	Event:   requests.EventRegister,
}

var jsonTestNotification, _ = json.Marshal(testNotification)

func setupDB() {
	mockNotifications := []requests.Notification{
		{Url: "testURL.com", Event: "REGISTER", Country: "NO"},
		{Url: "testURL.com", Event: "INVOKE", Country: "NO"},
		{Url: "testURL.com", Event: "INVOKE", Country: "SE"},
		{Url: "testURL.com", Event: "REGISTER", Country: ""},
	}
	for _, n := range mockNotifications {
		err := firebase.AddDocument[requests.Notification](n, firebase.NotificationCollection)
		if err != nil {
			log.Println("Error while trying to add notification to db: ", err.Error())
		}
	}
}

func TestMain(m *testing.M) {
	// Setup function
	log.Println("Setup for testing notifications")
	mock.InitForTesting()

	setupDB()

	// Run tests
	m.Run()

	defer func() {
		// Teardown function
		log.Println("Teardown for testing notifications")
		mock.TeardownAfterTesting()
	}()
}

func TestHandlerWithoutID(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statusCode int
	}{
		{
			name:       "NegativeTestNotificationHandlerWithoutID",
			method:     http.MethodOptions,
			statusCode: http.StatusNotImplemented,
		},
		{
			name:       "PositiveTestNotificationHandlerWithoutID",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a mock request
				req := httptest.NewRequest(tt.method, "/", nil)

				// Create a mock response recorder
				w := httptest.NewRecorder()

				// Call the handler
				HandlerWithoutID(w, req)

				// Check if the status code matches expected
				if w.Code != tt.statusCode {
					log.Println("Testing: ", tt.name)
					t.Errorf(
						"handler returned wrong status code: got %v want %v",
						w.Code, tt.statusCode,
					)
				}
			},
		)
	}
}

func Test_handleNotificationsGetRequest(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantedStatus int
	}{
		{
			name: "GetValidRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleNotificationsGetRequest(tt.args.w, tt.args.r)

				if tt.args.w.(*httptest.ResponseRecorder).Code != tt.wantedStatus {
					t.Errorf(
						"handleRegistrationsGetRequest() = %v, want %v",
						tt.args.w.(*httptest.ResponseRecorder).Code, tt.wantedStatus,
					)
				}
			},
		)
	}
}

func Test_handleNotificationsPostRequest(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantedStatus int
	}{
		{
			name: "PostValidRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonTestNotification)),
			},
			wantedStatus: http.StatusOK,
		},
		{
			name: "PostInvalidRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			wantedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleNotificationsPostRequest(tt.args.w, tt.args.r)

				switch tt.wantedStatus {
				case http.StatusOK:
					if tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusOK {
						t.Errorf(
							"handleNotificationsPostRequest() = %v, want %v",
							tt.args.w.(*httptest.ResponseRecorder).Code, http.StatusOK,
						)
					}
					if tt.args.w.(*httptest.ResponseRecorder).Body.String() == "" {
						t.Errorf(
							"handleNotificationsPostRequest() = %v, want an id in the response body",
							tt.args.w.(*httptest.ResponseRecorder).Body.String(),
						)
					}
				case http.StatusBadRequest:
					if tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusBadRequest {
						t.Errorf(
							"handleNotificationsPostRequest() = %v, want %v",
							tt.args.w.(*httptest.ResponseRecorder).Code, http.StatusBadRequest,
						)
					}
				}

			},
		)
	}
}

func Test_isValidEvent(t *testing.T) {
	type args struct {
		event string
	}
	tests := []struct {
		name    string
		args    args
		isValid bool
	}{
		{
			name: "IsValidEvent",
			args: args{
				event: requests.EventDelete,
			},
			isValid: true,
		},
		{
			name: "IsInvalidEvent",
			args: args{
				event: "INVALID_EVENT",
			},
			isValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := isValidEvent(tt.args.event); got != tt.isValid {
					t.Errorf("isValidEvent() = %v, want %v", got, tt.isValid)
				}
			},
		)
	}
}
