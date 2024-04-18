package notifications

import (
	"assignment-2/internal/constants"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWithID(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		method     string
		statusCode int
	}{
		{
			name:       "NegativeTestNotificationHandlerWithID",
			method:     http.MethodOptions,
			statusCode: http.StatusNotImplemented,
		},
		{
			name:       "NoIDTestNotificationHandlerWithID",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a mock request
				req := httptest.NewRequest(tt.method, "/", nil)

				// Create a mock response recorder
				w := httptest.NewRecorder()

				// Call the handler
				HandlerWithID(w, req)

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

func Test_handleNotificationsDeleteRequestWithID(t *testing.T) {
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
			name: "PositiveDeleteRequestWithID",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					constants.NotificationsPath+"?id="+getValidID(),
					nil,
				),
			},
			wantedStatus: http.StatusOK,
		},
		{
			name: "NegativeDeleteRequestWithID",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					constants.NotificationsPath+"?id=invalidID",
					nil,
				),
			},
			// Should maybe be bad request, but the db function returns internal server error
			wantedStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleNotificationsDeleteRequestWithID(tt.args.w, tt.args.r)

				if tt.args.w.(*httptest.ResponseRecorder).Code != tt.wantedStatus {
					t.Errorf(
						"handleNotificationsDeleteRequestWithID() = %v, want %v",
						tt.args.w.(*httptest.ResponseRecorder).Code, tt.wantedStatus,
					)
				}
			},
		)
	}
}

func Test_handleNotificationsGetRequestWithID(t *testing.T) {
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
			name: "PositiveGetRequestWithID",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					constants.NotificationsPath+"?id="+getValidID(),
					nil,
				),
			},
			wantedStatus: http.StatusOK,
		},
		{
			name: "NegativeGetRequestWithID",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					constants.NotificationsPath+"?id=invalidID",
					nil,
				),
			},
			// Should maybe be bad request, but the db function returns internal server error
			wantedStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleNotificationsGetRequestWithID(tt.args.w, tt.args.r)

				if tt.args.w.(*httptest.ResponseRecorder).Code != tt.wantedStatus {
					t.Errorf(
						"handleNotificationsGetRequestWithID() = %v, want %v",
						tt.args.w.(*httptest.ResponseRecorder).Code, tt.wantedStatus,
					)
				}
			},
		)
	}
}

func getValidID() string {
	// POST request to create a new registration, record the ID from the response

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonTestNotification))

	handleNotificationsPostRequest(w, r)

	// Get the ID from the response
	var notificationRes notificationResponse
	err := json.NewDecoder(w.Body).Decode(&notificationRes)
	if err != nil {
		log.Println("Error while decoding json: ", err.Error())
	}

	return notificationRes.Id
}
