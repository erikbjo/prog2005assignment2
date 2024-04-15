package dashboards

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/datasources/firebase"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/http/datatransfers/responses"
	utils2 "assignment-2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"
)

type dashboard struct {
	Country       string            `json:"country"`
	IsoCode       string            `json:"isoCode"`
	Features      dashboardFeatures `json:"features"`
	LastRetrieval time.Time         `json:"lastRetrieval"`
}

type dashboardFeatures struct {
	Temperature      float64                         `json:"temperature"`
	Precipitation    float64                         `json:"precipitation"`
	Capital          []string                        `json:"capital"`
	Coordinates      inhouse.Coordinates             `json:"coordinates"`
	Population       int                             `json:"population"`
	Area             float64                         `json:"area"`
	TargetCurrencies map[string]float64              `json:"targetCurrencies"`
	Currencies       map[string]responses.Currencies `json:"currencies"`
}

// Implemented methods for the endpoint
var implementedMethods = []string{
	http.MethodGet,
}

// Endpoint for managing dashboards
var dashboardsEndpoint = inhouse.Endpoint{
	Path:        constants.DashboardsPath,
	Methods:     implementedMethods,
	Description: "Endpoint for managing dashboards.",
}

// GetEndpointStructs returns the endpoint struct for the dashboards endpoint.
func GetEndpointStructs() []inhouse.Endpoint {
	return []inhouse.Endpoint{dashboardsEndpoint}
}

// HandlerWithID handles the /dashboard/v1/dashboards path.
// It currently only supports GET requests
func HandlerWithID(w http.ResponseWriter, r *http.Request) {
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleDashboardsGetRequest(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}

}

// handleDashboardsGetRequest handles the GET request for the /dashboard/v1/dashboards path.
// It is used to retrieve the populated dashboards.
func handleDashboardsGetRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils2.GetIDFromRequest(r)

	dashboardConfig, err := firebase.GetDocument[requests.DashboardConfig](
		id,
		firebase.DashboardCollection,
	)
	if err != nil {
		log.Println("Error while trying to display dashboard document: ", err.Error())
		http.Error(
			w,
			"Error while trying to display dashboard document",
			http.StatusInternalServerError,
		)
		return
	}

	var response dashboard
	response.Country = dashboardConfig.Country
	response.IsoCode = dashboardConfig.IsoCode

	log.Println("\n\tResponse: ", response)

	var features dashboardFeatures
	countryFeatures, err := getCountryData(dashboardConfig.IsoCode)
	if err != nil {
		log.Println("Error while trying to get country data: ", err.Error())
		http.Error(
			w,
			"Error while trying to get country data.",
			http.StatusInternalServerError,
		)
		return
	}

	// Merge the features
	// mergeFeatures(features, countryFeatures)

	features.Currencies = countryFeatures.Currencies
	features.Population = countryFeatures.Population
	features.Area = countryFeatures.Area
	features.Capital = countryFeatures.Capital
	features.Coordinates = countryFeatures.Coordinates

	// TODO: Get the weather data for the day, then taking mean
	meteoFeatures, err := getMeteoData(features.Coordinates)
	if err != nil {
		log.Println("Error while trying to get meteo data: ", err.Error())
		http.Error(
			w,
			"Error while trying to get meteo data.",
			http.StatusInternalServerError,
		)
		return
	}

	// Merge the features
	features.Temperature = meteoFeatures.Temperature
	features.Precipitation = meteoFeatures.Precipitation

	currencyFeatures, err := getCurrencyData(
		dashboardConfig.Features.TargetCurrencies,
		countryFeatures.Currencies,
	)
	if err != nil {
		log.Println("Error while trying to get currency rates: ", err.Error())
		http.Error(
			w,
			"Error while trying to get currency rates.",
			http.StatusInternalServerError,
		)
		return
	}

	// Merge the features
	features.TargetCurrencies = currencyFeatures.TargetCurrencies

	response.Features = features
	response.LastRetrieval = time.Now()

	log.Println("\n\tResponse: ", response)

	// Marshal the status object to JSON
	marshaled, err := json.MarshalIndent(
		response,
		"",
		"\t",
	)
	if err != nil {
		log.Println("Error during JSON encoding: " + err.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response
	_, err = w.Write(marshaled)
	if err != nil {
		log.Println("Failed to write response: " + err.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func getMeteoData(coordinates inhouse.Coordinates) (dashboardFeatures, error) {
	// Get the weather data from the meteo API
	r, err1 := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation&timezone=Europe%%2FBerlin&forecast_days=1",
			utils2.CurrentMeteoApi, coordinates.Latitude, coordinates.Longitude,
		),
		nil,
	)
	if err1 != nil {
		log.Println("Error in creating request:", err1.Error())
		return dashboardFeatures{}, fmt.Errorf("error in creating request")
	}

	log.Println("Request: ", r)

	r.Header.Add("content-type", "application/json")

	// Issue request
	res, err2 := utils2.Client.Do(r)
	if err2 != nil {
		log.Println("Error in response:", err2.Error())
		return dashboardFeatures{}, fmt.Errorf("error in response")
	}

	// Decode JSON
	var meteo responses.MeteoForecastResponse
	err3 := json.NewDecoder(res.Body).Decode(&meteo)
	if err3 != nil {
		log.Println("Error in decoding JSON:", err3.Error())
		return dashboardFeatures{}, fmt.Errorf("error in decoding JSON")
	}

	// features := dashboardFeatures{
	// 	Temperature:   meteo.Current.Temperature2M,
	// 	Precipitation: meteo.Current.Precipitation,
	// }
	// Gets the average of all hourly temperatures and rounds to 5 decimal points
	averageTemperature := float64(int(average(meteo.Hourly.Temperature2M)*100000)) / 100000
	averagePrecipitation := float64(int(average(meteo.Hourly.Precipitation)*100000)) / 100000

	fmt.Printf("request: %v\n", res)
	fmt.Printf("length of temp array: %d\n", len(meteo.Hourly.Temperature2M))
	features := dashboardFeatures{
		Temperature:   averageTemperature,
		Precipitation: averagePrecipitation,
	}

	return features, nil
}

func getCountryData(isoCode string) (dashboardFeatures, error) {
	// Get the country data from the restcountries API
	r, err1 := http.NewRequest(
		http.MethodGet,
		utils2.CurrentRestCountriesApi+"alpha/"+isoCode+"?fields=name,cca2,currencies,capital,latlng,area,population",
		nil,
	)
	if err1 != nil {
		log.Println("Error in creating request:", err1.Error())
		return dashboardFeatures{}, fmt.Errorf("error in creating request")
	}

	r.Header.Add("content-type", "application/json")

	// Issue request
	res, err2 := utils2.Client.Do(r)
	if err2 != nil {
		log.Println("Error in response:", err2.Error())
		return dashboardFeatures{}, fmt.Errorf("error in response")
	}

	// Decode JSON
	var country responses.ResponseFromRestcountries
	err3 := json.NewDecoder(res.Body).Decode(&country)
	if err3 != nil {
		log.Println("Error in decoding JSON:", err3.Error())
		return dashboardFeatures{}, fmt.Errorf("error in decoding JSON")
	}

	if country.Name.Common == "" {
		log.Println("Country not found")
		return dashboardFeatures{}, fmt.Errorf("country not found")
	}

	lat := country.Latlng[0]
	lng := country.Latlng[1]

	features := dashboardFeatures{
		Capital: country.Capital,
		Coordinates: inhouse.Coordinates{
			Latitude:  lat,
			Longitude: lng,
		},
		Population: country.Population,
		Area:       country.Area,
		Currencies: country.Currencies,
	}

	return features, nil
}

func getCurrencyData(
	targetCurrencies []string,
	currencies map[string]responses.Currencies,
) (dashboardFeatures, error) {
	featuresFromCurrency := dashboardFeatures{
		TargetCurrencies: make(map[string]float64),
	}
	// Task specifies to take the first currency where multiple currencies are available
	// Therefore, we will only take the first currency from the map
	if len(currencies) == 0 {
		log.Println("No currencies found")
		return dashboardFeatures{}, fmt.Errorf("no currencies found")
	}

	var exchangeCurrency string
	for key := range currencies {
		exchangeCurrency = key
		break
	}

	// Get the exchange rates from the currency API
	r, err1 := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s",
			utils2.CurrentCurrencyApi, exchangeCurrency,
		),
		nil,
	)
	if err1 != nil {
		log.Println("Error in creating request:", err1.Error())
		return dashboardFeatures{}, fmt.Errorf("error in creating request")
	}

	r.Header.Add("content-type", "application/json")

	// Issue request
	res, err2 := utils2.Client.Do(r)
	if err2 != nil {
		log.Println("Error in response:", err2.Error())
		return dashboardFeatures{}, fmt.Errorf("error in response")
	}

	// Decode JSON
	var response responses.ResponseFromCurrency
	err3 := json.NewDecoder(res.Body).Decode(&response)
	if err3 != nil {
		log.Println("Error in decoding JSON:", err3.Error())
		return dashboardFeatures{}, fmt.Errorf("error in decoding JSON")
	}

	// Get the exchange rates for the target currencies
	for _, targetCurrency := range targetCurrencies {
		if _, ok := response.Rates[targetCurrency]; !ok {
			log.Println("Exchange rate not found for currency: ", targetCurrency)
			// Not returning error, just setting the rate to 0
			featuresFromCurrency.TargetCurrencies["targetCurrency"] = 0
		} else {
			featuresFromCurrency.TargetCurrencies[targetCurrency] = response.Rates[targetCurrency]
		}
	}

	return featuresFromCurrency, nil
}

func mergeFeatures(a, b dashboardFeatures) {
	ra := reflect.ValueOf(&a).Elem()
	rb := reflect.ValueOf(&b).Elem()

	numFields := ra.NumField()

	for i := 0; i < numFields; i++ {
		log.Println(
			"\n\tChecking field: ", ra.Type().Field(i).Name+
				" with type: ", ra.Field(i).Kind(),
		)
		fieldA := ra.Field(i)
		fieldB := rb.Field(i)

		switch fieldA.Kind() {
		// case reflect.Float64, reflect.Kind():
		case reflect.Float64, reflect.Slice, reflect.Struct, reflect.Int, reflect.Map:
			// case reflect.Ptr:

			if fieldA.IsNil() {
				fieldA.Set(fieldB)
			}
		default:
			log.Println("Unsupported type when merging features: ", fieldA.Kind())
		}
	}
}

func average(elements []float64) float64 {
	var sum float64
	if len(elements) == 0 {
		return 0
	}
	for _, element := range elements {
		sum += element
	}
	return sum / float64(len(elements))
}
