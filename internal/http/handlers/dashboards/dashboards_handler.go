package dashboards

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/db"
	"assignment-2/internal/http/datatransfers/inhouse"
	"assignment-2/internal/http/datatransfers/requests"
	"assignment-2/internal/http/datatransfers/responses"
	utils2 "assignment-2/internal/utils"
	"dario.cat/mergo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// dashboard is the struct for the response object
type dashboard struct {
	Country       string            `json:"country"`
	IsoCode       string            `json:"isoCode"`
	Features      dashboardFeatures `json:"features"`
	LastRetrieval time.Time         `json:"lastRetrieval"`
}

// dashboardFeatures is the struct for the features of the dashboard
type dashboardFeatures struct {
	Temperature      *float64             `json:"temperature,omitempty"`
	Precipitation    *float64             `json:"precipitation,omitempty"`
	Capital          *string              `json:"capital,omitempty"`
	Coordinates      *inhouse.Coordinates `json:"coordinates,omitempty"`
	Population       *int                 `json:"population,omitempty"`
	Area             *float64             `json:"area,omitempty"`
	TargetCurrencies map[string]float64   `json:"targetCurrencies,omitempty"`
	Currency         responses.Currency   `json:"currency"`
}

// Implemented methods for the endpoint
var implementedMethods = []string{
	http.MethodGet,
}

// Endpoint for managing dashboards
var dashboardsEndpoint = inhouse.Endpoint{
	Path:        constants.DashboardsPath + "{id}",
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
	w.Header().Set("content-type", "application/json")
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

	dashboardConfig, err := db.GetDocument[requests.DashboardConfig](
		id,
		db.DashboardCollection,
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

	// Create the response object and assign the country and iso code
	var response dashboard
	response.Country = dashboardConfig.Country
	response.IsoCode = dashboardConfig.IsoCode

	// Get the features for the dashboard
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
	err = mergo.Merge(&features, countryFeatures, mergo.WithOverride, mergo.WithoutDereference)
	if err != nil {
		log.Println("Error while trying to merge country features: ", err.Error())
		http.Error(
			w,
			"Error while trying to merge country features.",
			http.StatusInternalServerError,
		)
		return
	}

	// Get the meteo features
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
	err = mergo.Merge(&features, meteoFeatures, mergo.WithOverride, mergo.WithoutDereference)
	if err != nil {
		log.Println("Error while trying to merge meteo features: ", err.Error())
		http.Error(
			w,
			"Error while trying to merge meteo features.",
			http.StatusInternalServerError,
		)
		return
	}

	// Get the currency features
	currencyFeatures, err := getCurrencyData(
		dashboardConfig.Features.TargetCurrencies,
		countryFeatures.Currency,
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
	err = mergo.Merge(&features, currencyFeatures, mergo.WithOverride, mergo.WithoutDereference)
	if err != nil {
		log.Println("Error while trying to merge currency features: ", err.Error())
		http.Error(
			w,
			"Error while trying to merge currency features.",
			http.StatusInternalServerError,
		)
		return
	}

	// Assign the features to the response
	response.Features = features
	response.LastRetrieval = time.Now()

	// Filter the response by the config
	filteredResponse, err := filterDashboardByConfig(response, dashboardConfig)
	if err != nil {
		log.Println("Error while trying to filter dashboard features: ", err.Error())
		http.Error(
			w,
			"Error while trying to filter dashboard features.",
			http.StatusInternalServerError,
		)
		return
	}

	// Marshal the status object to JSON
	marshaled, err := json.MarshalIndent(
		filteredResponse,
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

// getMeteoData gets the meteo data for the given coordinates.
// This data includes the mean temperature and precipitation.
func getMeteoData(coordinates *inhouse.Coordinates) (dashboardFeatures, error) {
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

	// Gets the average of all hourly temperatures and rounds to 5 decimal points
	averageTemperature := float64(int(average(meteo.Hourly.Temperature2M)*100000)) / 100000
	averagePrecipitation := float64(int(average(meteo.Hourly.Precipitation)*100000)) / 100000

	features := dashboardFeatures{
		Temperature:   &averageTemperature,
		Precipitation: &averagePrecipitation,
	}

	return features, nil
}

// getCountryData gets the country data for the given ISO code. This data includes the capital, coordinates, population,
// area, and currency.
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

	log.Println("Request url: ", r.URL)

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
	// Task specifies to take the first capital where multiple capitals are available
	capital := country.Capital[0]
	// Same goes for currency
	var currency responses.Currency
	for key, value := range country.Currencies {
		currency = value
		currency.Code = key
		break
	}

	features := dashboardFeatures{
		Capital: &capital,
		Coordinates: &inhouse.Coordinates{
			Latitude:  lat,
			Longitude: lng,
		},
		Population: &country.Population,
		Area:       &country.Area,
		Currency:   currency,
	}

	return features, nil
}

// getCurrencyData gets the currency data for the given target currencies. This data includes the exchange rates.
func getCurrencyData(
	targetCurrencies []string,
	exchangeCurrency responses.Currency,
) (dashboardFeatures, error) {
	featuresFromCurrency := dashboardFeatures{
		TargetCurrencies: make(map[string]float64),
	}
	// Get the exchange rates from the currency API
	r, err1 := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s%s",
			utils2.CurrentCurrencyApi, exchangeCurrency.Code,
		),
		nil,
	)
	if err1 != nil {
		log.Println("Error in creating request:", err1.Error())
		return dashboardFeatures{}, fmt.Errorf("error in creating request")
	}

	r.Header.Add("content-type", "application/json")
	log.Println("Request: ", r)

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

// average calculates the mean of a slice of float64 elements.
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

// filterDashboardByConfig filters the dashboard features by the given config.
func filterDashboardByConfig(oldDashboard dashboard, config requests.DashboardConfig) (
	dashboard,
	error,
) {
	if config.Country != oldDashboard.Country {
		return dashboard{}, fmt.Errorf("country does not match")
	}
	// Returns a new dashboard with the features filtered by the config
	newDashboard := dashboard{
		Country:       oldDashboard.Country,
		IsoCode:       oldDashboard.IsoCode,
		LastRetrieval: oldDashboard.LastRetrieval,
	}

	if config.Features.Temperature {
		newDashboard.Features.Temperature = oldDashboard.Features.Temperature
	}
	if config.Features.Precipitation {
		newDashboard.Features.Precipitation = oldDashboard.Features.Precipitation
	}
	if config.Features.Capital {
		newDashboard.Features.Capital = oldDashboard.Features.Capital
	}
	if config.Features.Coordinates {
		newDashboard.Features.Coordinates = oldDashboard.Features.Coordinates
	}
	if config.Features.Population {
		newDashboard.Features.Population = oldDashboard.Features.Population
	}
	if config.Features.Area {
		newDashboard.Features.Area = oldDashboard.Features.Area
	}

	// Translocate the target currencies
	newDashboard.Features.TargetCurrencies = make(map[string]float64)
	for key, value := range oldDashboard.Features.TargetCurrencies {
		newDashboard.Features.TargetCurrencies[key] = value
	}

	newDashboard.Features.Currency = oldDashboard.Features.Currency

	return newDashboard, nil
}
