package registrations

import (
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/mock"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testRegistration = requests.DashboardConfig{
	Country: "Norway",
	IsoCode: "NO",
	Features: requests.ConfigFeatures{
		Temperature:      true,
		Precipitation:    true,
		Capital:          true,
		Coordinates:      true,
		Population:       true,
		Area:             true,
		TargetCurrencies: []string{"USD", "EUR"},
	},
}

var jsonTestRegistration, _ = json.Marshal(testRegistration)

func TestMain(m *testing.M) {
	// Setup function
	log.Println("Setup for testing")
	mock.InitForTesting()

	// Run tests
	m.Run()

	// Teardown function
	log.Println("Teardown for testing")
	mock.TeardownAfterTesting()
}

func TestHandlerWithoutID(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statusCode int
	}{
		{
			name:       "NegativeTestRegistrationHandlerWithoutID",
			method:     http.MethodOptions,
			statusCode: http.StatusNotImplemented,
		},
		{
			name:       "PositiveTestRegistrationHandlerWithoutID",
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

func Test_handleRegistrationsGetRequest(t *testing.T) {
	// Since we are using a test database, we can't test the actual functionality of this and next function
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
				handleRegistrationsGetRequest(tt.args.w, tt.args.r)

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

func Test_handleRegistrationsHeadRequest(t *testing.T) {
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
			name: "HeadValidRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodHead, "/", nil),
			},
			wantedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleRegistrationsHeadRequest(tt.args.w, tt.args.r)

				if tt.args.w.(*httptest.ResponseRecorder).Code != tt.wantedStatus {
					t.Errorf(
						"handleRegistrationsHeadRequest() = %v, want %v",
						tt.args.w.(*httptest.ResponseRecorder).Code, tt.wantedStatus,
					)
				}
			},
		)
	}
}

func Test_handleRegistrationsPostRequest(t *testing.T) {
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
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonTestRegistration)),
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
				handleRegistrationsPostRequest(tt.args.w, tt.args.r)

				switch tt.wantedStatus {
				case http.StatusOK:
					if tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusOK {
						t.Errorf(
							"handleRegistrationsPostRequest() = %v, want %v",
							tt.args.w.(*httptest.ResponseRecorder).Code, http.StatusOK,
						)
					}
					if tt.args.w.(*httptest.ResponseRecorder).Body.String() == "" {
						t.Errorf(
							"handleRegistrationsPostRequest() = %v, want an id in the response body",
							tt.args.w.(*httptest.ResponseRecorder).Body.String(),
						)
					}
				case http.StatusBadRequest:
					if tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusBadRequest {
						t.Errorf(
							"handleRegistrationsPostRequest() = %v, want %v",
							tt.args.w.(*httptest.ResponseRecorder).Code, http.StatusBadRequest,
						)
					}
				}
			},
		)
	}
}
