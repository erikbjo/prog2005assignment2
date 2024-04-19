package status

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/mock"
	"assignment-2/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup function
	log.Println("Setup for testing status")
	mock.InitForTesting()

	// Run tests
	m.Run()

	// Teardown function
	defer func() {
		log.Println("Teardown for testing status")
		mock.TeardownAfterTesting()
	}()

}

// TestStatusHandler tests the Handler function, which handles requests for /status
// It tests the GET method for the /status path
func TestStatusHandler(t *testing.T) {
	// Create tests with different HTTP methods and expected status codes
	tests := []struct {
		name       string
		method     string
		statusCode int
	}{
		{
			name:       "NegativeTestStatusHandler",
			method:     http.MethodOptions,
			statusCode: http.StatusNotImplemented,
		},
		{
			name:       "PositiveTestStatusHandler",
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
				Handler(w, req)

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

func Test_getStatusCode(t *testing.T) {
	type args struct {
		url string
		w   http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test_getStatusCodeCurrentRestCountriesApi",
			args: args{
				url: utils.CurrentRestCountriesApi,
				w:   httptest.NewRecorder(),
			},
			want: http.StatusOK,
		},
		{
			name: "Test_getStatusCodeCurrentMeteoApi",
			args: args{
				url: utils.CurrentMeteoApi,
				w:   httptest.NewRecorder(),
			},
			want: http.StatusOK,
		},
		{
			name: "Test_getStatusCodeCurrentCurrencyApi",
			args: args{
				url: utils.CurrentCurrencyApi,
				w:   httptest.NewRecorder(),
			},
			want: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := getStatusCode(tt.args.url, tt.args.w); got != tt.want {
					t.Errorf("getStatusCode() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_handleStatusGetRequest(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_handleStatusGetRequest",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, constants.StatusPath, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleStatusGetRequest(tt.args.w, tt.args.r)

				// Test the response
				if tt.args.w.(*httptest.ResponseRecorder).Code != http.StatusOK {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						tt.args.w.(*httptest.ResponseRecorder).Code,
						http.StatusOK,
					)
				}

				if tt.args.w.(*httptest.ResponseRecorder).Body.String() == "" {
					t.Errorf(
						"handleStatusGetRequest() = %v, want not empty",
						tt.args.w.(*httptest.ResponseRecorder).Body.String(),
					)
				}

				// Make body a status object
				var status status
				err := json.Unmarshal(tt.args.w.(*httptest.ResponseRecorder).Body.Bytes(), &status)
				if err != nil {
					t.Errorf("handleStatusGetRequest() = %v", err)
				}

				// Check if the status object is correct
				if status.Version != constants.Version {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						status.Version,
						constants.Version,
					)
				}

				// Check that uptime is an int
				if reflect.TypeOf(status.Uptime) != reflect.TypeOf(int(0)) {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						reflect.TypeOf(status.Uptime),
						reflect.Int,
					)
				}

				if status.CountriesAPI != http.StatusOK {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						status.CountriesAPI,
						http.StatusOK,
					)
				}

				if status.MeteoAPI != http.StatusOK {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						status.MeteoAPI,
						http.StatusOK,
					)
				}

				if status.CurrencyAPI != http.StatusOK {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						status.CurrencyAPI,
						http.StatusOK,
					)
				}

				if status.DashboardDB != http.StatusOK {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						status.DashboardDB,
						http.StatusOK,
					)
				}

				if status.NotificationDB != http.StatusOK {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						status.NotificationDB,
						http.StatusOK,
					)
				}

				// TODO: Implement counting of webhooks
				if status.Webhooks != http.StatusNotImplemented {
					t.Errorf(
						"handleStatusGetRequest() = %v, want %v",
						status.Webhooks,
						http.StatusNotImplemented,
					)
				}

			},
		)
	}
}
