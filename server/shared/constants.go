package shared

// DefaultPath Default path for the server
const DefaultPath = "/"

// Version Version of the API
const Version = "v1"

// DashboardPath Path for the dashboard
const DashboardPath = "/dashboard/" + Version

// RegistrationsPath Path for the registrations
const RegistrationsPath = DashboardPath + "/registrations/"

// DashboardsPath Path for the dashboards
const DashboardsPath = DashboardPath + "/dashboards/"

// NotificationsPath Path for the notifications
const NotificationsPath = DashboardPath + "/notifications/"

// StatusPath Path for the status
const StatusPath = DashboardPath + "/status/"

// RestCountriesApi Christopher's RestCountries API
const RestCountriesApi = "http://129.241.150.113:8080/v3.1"

// TestRestCountriesApi Local stub for RestCountries API
const TestRestCountriesApi = "http://localhost:8080/test/restcountries"

// CurrencyApi Christopher's Currency API
const CurrencyApi = "http://129.241.150.113:9090/currency/"

// TestCurrencyApi Local stub for Currency API
const TestCurrencyApi = "http://localhost:8080/test/currency"

/* https://open-meteo.com/en/features#available-apis */
