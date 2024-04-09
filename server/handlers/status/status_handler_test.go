package status

import (
	"assignment-2/server/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestStatusHandler tests the Handler function, which handles requests for /status
// It tests the GET method for the /status path
func TestStatusHandler(t *testing.T) {
	// Use stubs for testing
	handlers.SetStubsForTesting()

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
					t.Errorf(
						"handler returned wrong status code: got %v want %v",
						w.Code, tt.statusCode,
					)
				}
			},
		)
	}
}

// TODO: Implement the stubs to test
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
		// TODO: Add test cases.
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

// The status codes are covered by the Test_getStatusCode function
// TODO: Implement the stubs to test
func Test_handleStatusGetRequest(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handleStatusGetRequest(tt.args.w, tt.args.r)
			},
		)
	}
}
