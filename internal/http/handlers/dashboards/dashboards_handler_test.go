package dashboards

import (
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/http/datatransfers/responses"
	"assignment-2/internal/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

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

func TestHandlerWithID(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statusCode int
	}{
		{
			name:       "NegativeTestHandlerWithID",
			method:     http.MethodOptions,
			statusCode: http.StatusNotImplemented,
		},
		{
			name:       "PositiveTestHandlerWithID",
			method:     http.MethodGet,
			statusCode: http.StatusInternalServerError,
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

func Test_average(t *testing.T) {
	type args struct {
		elements []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test_average",
			args: args{
				elements: []float64{1, 2, 3, 4, 5},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := average(tt.args.elements); got != tt.want {
					t.Errorf("average() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_filterDashboardByConfig(t *testing.T) {
	type args struct {
		oldDashboard dashboard
		config       requests.DashboardConfig
	}
	tests := []struct {
		name    string
		args    args
		want    dashboard
		wantErr bool
	}{
		{
			name: "Test_filterDashboardByConfigAllFalse",
			args: args{
				oldDashboard: dashboard{},
				config: requests.DashboardConfig{
					ID:      "",
					Country: "",
					IsoCode: "",
					Features: requests.ConfigFeatures{
						Temperature:      false,
						Precipitation:    false,
						Capital:          false,
						Coordinates:      false,
						Population:       false,
						Area:             false,
						TargetCurrencies: nil,
					},
					LastChange: time.Time{},
				},
			},
			want:    dashboard{},
			wantErr: false,
		},
		{
			name: "Test_filterDashboardByConfigAllTrue",
			args: args{
				oldDashboard: dashboard{
					Country: "Test",
					IsoCode: "Test",
					Features: dashboardFeatures{
						Temperature:      new(float64),
						Precipitation:    new(float64),
						Capital:          new(string),
						Coordinates:      &inhouse.Coordinates{},
						Population:       new(int),
						Area:             new(float64),
						TargetCurrencies: map[string]float64{},
						Currency:         responses.Currency{},
					},
					LastRetrieval: time.Time{},
				},
				config: requests.DashboardConfig{
					ID:      "Test",
					Country: "Test",
					IsoCode: "Test",
					Features: requests.ConfigFeatures{
						Temperature:      true,
						Precipitation:    true,
						Capital:          true,
						Coordinates:      true,
						Population:       false,
						Area:             false,
						TargetCurrencies: nil,
					},
					LastChange: time.Time{},
				},
			},
			want: dashboard{
				Country: "Test",
				IsoCode: "Test",
				Features: dashboardFeatures{
					Temperature:      new(float64),
					Precipitation:    new(float64),
					Capital:          new(string),
					Coordinates:      &inhouse.Coordinates{},
					Population:       nil,
					Area:             nil,
					TargetCurrencies: nil,
					Currency:         responses.Currency{},
				},
				LastRetrieval: time.Time{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := filterDashboardByConfig(tt.args.oldDashboard, tt.args.config)
				if (err != nil) != tt.wantErr {
					t.Errorf("filterDashboardByConfig() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// Check if the country is correct
				if got.Country != tt.want.Country {
					t.Errorf(
						"wrong country: got %v want %v",
						got.Country,
						tt.want.Country,
					)
				}

				// Check if the iso code is correct
				if got.IsoCode != tt.want.IsoCode {
					t.Errorf(
						"wrong iso code: got %v want %v",
						got.IsoCode,
						tt.want.IsoCode,
					)
				}

				// Check if the temperature is correct
				if reflect.TypeOf(got.Features.Temperature).Kind() != reflect.Ptr && reflect.TypeOf(got.Features.Temperature).Elem().Kind() != reflect.Float64 {
					t.Errorf(
						"wrong type for temperature: got %v want %v",
						reflect.TypeOf(got.Features.Temperature),
						reflect.TypeOf(tt.want.Features.Temperature),
					)
				}

				// Check if the precipitation is correct
				if reflect.TypeOf(got.Features.Precipitation).Kind() != reflect.Ptr && reflect.TypeOf(got.Features.Precipitation).Elem().Kind() != reflect.Float64 {
					t.Errorf(
						"wrong type for precipitation: got %v want %v",
						reflect.TypeOf(got.Features.Precipitation),
						reflect.TypeOf(tt.want.Features.Precipitation),
					)
				}

				// Check if the capital is correct
				if reflect.TypeOf(got.Features.Capital).Kind() != reflect.Ptr && reflect.TypeOf(got.Features.Capital).Elem().Kind() != reflect.String {
					t.Errorf(
						"wrong type for capital: got %v want %v",
						reflect.TypeOf(got.Features.Capital),
						reflect.TypeOf(tt.want.Features.Capital),
					)
				}

				// Check if the coordinates are correct
				if reflect.TypeOf(got.Features.Coordinates).Kind() != reflect.Ptr && reflect.TypeOf(got.Features.Coordinates).Elem().Kind() != reflect.Struct {
					t.Errorf(
						"wrong type for coordinates: got %v want %v",
						reflect.TypeOf(got.Features.Coordinates),
						reflect.TypeOf(tt.want.Features.Coordinates),
					)
				}

				// Check if the population is correct
				if reflect.TypeOf(got.Features.Population).Kind() != reflect.Ptr && reflect.TypeOf(got.Features.Population).Elem().Kind() != reflect.Int {
					t.Errorf(
						"wrong type for population: got %v want %v",
						reflect.TypeOf(got.Features.Population),
						reflect.TypeOf(tt.want.Features.Population),
					)
				}

				// Check if the area is correct
				if reflect.TypeOf(got.Features.Area).Kind() != reflect.Ptr && reflect.TypeOf(got.Features.Area).Elem().Kind() != reflect.Float64 {
					t.Errorf(
						"wrong type for area: got %v want %v",
						reflect.TypeOf(got.Features.Area),
						reflect.TypeOf(tt.want.Features.Area),
					)
				}

				// Check if the target currencies are correct
				if reflect.TypeOf(got.Features.TargetCurrencies).Kind() != reflect.Map {
					t.Errorf(
						"wrong type for target currencies: got %v want %v",
						reflect.TypeOf(got.Features.TargetCurrencies),
						reflect.TypeOf(tt.want.Features.TargetCurrencies),
					)
				}

				// Check if the currency is correct
				if reflect.TypeOf(got.Features.Currency).Kind() != reflect.Struct {
					t.Errorf(
						"wrong type for currency: got %v want %v",
						reflect.TypeOf(got.Features.Currency),
						reflect.TypeOf(tt.want.Features.Currency),
					)
				}

				// Check if the last retrieval is correct
				if got.LastRetrieval != tt.want.LastRetrieval {
					t.Errorf(
						"wrong last retrieval: got %v want %v",
						got.LastRetrieval,
						tt.want.LastRetrieval,
					)
				}
			},
		)
	}
}

func Test_getCountryData(t *testing.T) {
	type args struct {
		isoCode string
	}
	tests := []struct {
		name    string
		args    args
		want    dashboardFeatures
		wantErr bool
	}{
		{
			name: "Test_getCountryData",
			args: args{
				isoCode: "NO",
			},
			want: dashboardFeatures{
				Capital:    new(string),
				Population: new(int),
				Area:       new(float64),
				Currency: responses.Currency{
					Name:   "Norwegian krone",
					Symbol: "kr",
					Code:   "NOK",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getCountryData(tt.args.isoCode)
				if (err != nil) != tt.wantErr {
					t.Errorf("getCountryData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// Capital is a pointer to a string
				if reflect.TypeOf(got.Capital).Kind() != reflect.Ptr && reflect.TypeOf(got.Capital).Elem().Kind() != reflect.String {
					t.Errorf(
						"wrong type for capital: got %v want %v",
						reflect.TypeOf(got.Capital),
						reflect.TypeOf(tt.want.Capital),
					)
				}

				// Population is a pointer to an int
				if reflect.TypeOf(got.Population).Kind() != reflect.Ptr && reflect.TypeOf(got.Population).Elem().Kind() != reflect.Int {
					t.Errorf(
						"wrong type for population: got %v want %v",
						reflect.TypeOf(got.Population),
						reflect.TypeOf(tt.want.Population),
					)
				}

				// Area is a pointer to a float64
				if reflect.TypeOf(got.Area).Kind() != reflect.Ptr && reflect.TypeOf(got.Area).Elem().Kind() != reflect.Float64 {
					t.Errorf(
						"wrong type for area: got %v want %v",
						reflect.TypeOf(got.Area),
						reflect.TypeOf(tt.want.Area),
					)
				}

				// Currency is a struct
				if reflect.TypeOf(got.Currency).Kind() != reflect.Struct {
					t.Errorf(
						"wrong type for currency: got %v want %v",
						reflect.TypeOf(got.Currency),
						reflect.TypeOf(tt.want.Currency),
					)
				}

				// Check code of currency
				if got.Currency.Code != tt.want.Currency.Code {
					t.Errorf(
						"wrong code for currency: got %v want %v",
						got.Currency.Code,
						tt.want.Currency.Code,
					)
				}
			},
		)
	}
}

func Test_getCurrencyData(t *testing.T) {
	type args struct {
		targetCurrencies []string
		exchangeCurrency responses.Currency
	}
	tests := []struct {
		name    string
		args    args
		want    dashboardFeatures
		wantErr bool
	}{
		{
			name: "Test_getCurrencyData",
			args: args{
				targetCurrencies: []string{"USD", "EUR"},
				exchangeCurrency: responses.Currency{
					Name:   "Norwegian krone",
					Symbol: "kr",
					Code:   "NOK",
				},
			},
			want: dashboardFeatures{
				TargetCurrencies: map[string]float64{
					"USD": 0.093687, // Magic number from json file
					"EUR": 0.086289, // Magic number from json file
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getCurrencyData(tt.args.targetCurrencies, tt.args.exchangeCurrency)
				if (err != nil) != tt.wantErr {
					t.Errorf("getCurrencyData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// TargetCurrencies is a map
				if reflect.TypeOf(got.TargetCurrencies).Kind() != reflect.Map {
					t.Errorf(
						"wrong type for targetCurrencies: got %v want %v",
						reflect.TypeOf(got.TargetCurrencies),
						reflect.TypeOf(tt.want.TargetCurrencies),
					)
				}

				// Check if the target currencies are correct
				for key, value := range got.TargetCurrencies {
					if value != tt.want.TargetCurrencies[key] {
						t.Errorf(
							"wrong value for target currency: got %v want %v",
							value,
							tt.want.TargetCurrencies[key],
						)
					}
				}

				if len(got.TargetCurrencies) != len(tt.want.TargetCurrencies) {
					t.Errorf(
						"wrong number of target currencies: got %v want %v",
						len(got.TargetCurrencies),
						len(tt.want.TargetCurrencies),
					)
				}
			},
		)
	}
}

func Test_getMeteoData(t *testing.T) {
	type args struct {
		coordinates *inhouse.Coordinates
	}
	tests := []struct {
		name    string
		args    args
		want    dashboardFeatures
		wantErr bool
	}{
		{
			name: "Test_getMeteoData",
			args: args{
				coordinates: &inhouse.Coordinates{
					// Coordinates does not matter for the test
					Latitude:  60,
					Longitude: 11,
				},
			},
			want: dashboardFeatures{
				Temperature:   new(float64),
				Precipitation: new(float64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getMeteoData(tt.args.coordinates)
				if (err != nil) != tt.wantErr {
					t.Errorf("getMeteoData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// Reflect that precipitation is a pointer to a float64
				if reflect.TypeOf(got.Precipitation).Kind() != reflect.Ptr && reflect.TypeOf(got.Precipitation).Elem().Kind() != reflect.Float64 {
					t.Errorf(
						"wrong type for precipitation: got %v want %v",
						reflect.TypeOf(got.Precipitation),
						reflect.TypeOf(tt.want.Precipitation),
					)
				}
			},
		)
	}
}

// Since all other functions are tested, we skip handler function,
// as the untested code relies on the external functions
