package registrations

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
	tests := []struct {
		name       string
		method     string
		statusCode int
	}{
		{
			name:       "NegativeTestRegistrationHandlerWithID",
			method:     http.MethodOptions,
			statusCode: http.StatusNotImplemented,
		},
		{
			name:       "NoIDTestRegistrationHandlerWithID",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
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

func Test_handleRegistrationsDeleteRequestWithID(t *testing.T) {
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
					constants.RegistrationsPath+"?id="+getValidID(),
					nil,
				),
			},
			wantedStatus: http.StatusNoContent,
		},
		{
			name: "NegativeDeleteRequestWithID",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					constants.RegistrationsPath+"?id=",
					nil,
				),
			},
			// Should maybe be bad request, but the db function returns internal server error
			wantedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleRegistrationsDeleteRequestWithID(tt.args.w, tt.args.r)

				if tt.args.w.(*httptest.ResponseRecorder).Code != tt.wantedStatus {
					t.Errorf(
						"handleRegistrationsDeleteRequestWithID() = %v, want %v",
						tt.args.w.(*httptest.ResponseRecorder).Code, tt.wantedStatus,
					)
				}
			},
		)
	}
}

func Test_handleRegistrationsGetRequestWithID(t *testing.T) {
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
					constants.RegistrationsPath+"?id="+getValidID(),
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
					constants.RegistrationsPath+"?id=",
					nil,
				),
			},
			// Should maybe be bad request, but the db function returns internal server error
			wantedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleRegistrationsGetRequestWithID(tt.args.w, tt.args.r)

				if tt.args.w.(*httptest.ResponseRecorder).Code != tt.wantedStatus {
					t.Errorf(
						"handleRegistrationsGetRequestWithID() = %v, want %v",
						tt.args.w.(*httptest.ResponseRecorder).Code, tt.wantedStatus,
					)
				}
			},
		)
	}
}

func Test_handleRegistrationsPutRequestWithID(t *testing.T) {
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
			name: "PositivePutRequestWithID",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					constants.RegistrationsPath+"?id="+getValidID(),
					bytes.NewBuffer(jsonTestRegistration),
				),
			},
			wantedStatus: http.StatusNoContent,
		},
		{
			name: "NegativePutRequestWithID",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					constants.RegistrationsPath+"?id=invalidID",
					bytes.NewBuffer(jsonTestRegistration),
				),
			},
			wantedStatus: http.StatusInternalServerError,
		},
		{
			name: "NegativePutRequestWithBadBody",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					constants.RegistrationsPath+"?id="+getValidID(),
					bytes.NewBuffer([]byte("invalid json")),
				),
			},
			wantedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleRegistrationsPutRequestWithID(tt.args.w, tt.args.r)

				if tt.args.w.(*httptest.ResponseRecorder).Code != tt.wantedStatus {
					t.Errorf(
						"handleRegistrationsPutRequestWithID() = %v, want %v",
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
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonTestRegistration))

	handleRegistrationsPostRequest(w, r)

	// Get the ID from the response
	var registrationRes registrationResponse
	err := json.NewDecoder(w.Body).Decode(&registrationRes)
	if err != nil {
		log.Println("Error while decoding json: ", err.Error())
	}

	return registrationRes.ID
}
