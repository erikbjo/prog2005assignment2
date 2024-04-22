package constants

var (
	ErrJsonMarshal   = "error marshalling JSON"
	ErrJsonWrite     = "error writing JSON"
	ErrJsonEncode    = "error encoding JSON"
	ErrJsonDecode    = "error decoding JSON"
	ErrJsonRead      = "error reading JSON"
	ErrJsonUnmarshal = "error unmarshalling JSON"
	ErrJsonParse     = "error parsing JSON"
	ErrJsonInvalid   = "invalid JSON"

	ErrDBCount       = "error getting count from database"
	ErrDBRead        = "error reading from database"
	ErrDBWrite       = "error writing to database"
	ErrDBAddDoc      = "error adding document to database"
	ErrDBUpdateDoc   = "error updating document in database"
	ErrDBDeleteDoc   = "error deleting document from database"
	ErrDBGetDoc      = "error getting document from database"
	ErrDBClose       = "error closing database"
	ErrDBDocNotFound = "document not found in collection"
	ErrDBNoDocs      = "no documents found in collection"

	ErrFirestoreClient        = "error creating firestore client"
	ErrFirestoreClose         = "error closing firestore client"
	ErrFirestoreApp           = "error while creating firestore app"
	ErrFirestoreEmulatorEnv   = "error while setting up firestore emulator env variable"
	ErrFirestoreClientNotInit = "firestore client not initialized"

	ErrExternalResponse = "error getting response from external service"
	ErrExternalRequest  = "error making request to external service"

	ErrWriteResponse = "error writing response"

	ErrIDFromRequest = "error getting ID from request"
	ErrIDRequired    = "ID is required for endpoints with path ending in {id}."
	ErrIDInvalid     = "invalid ID provided"
	ErrIDNotProvided = "no ID provided"

	ErrDataNotMatchingTargetStruct = "data does not match target struct"

	ErrNotificationsInvalidType  = "invalid event type provided"
	ErrNotificationsGetDocFromDB = "error getting notification document from database"

	ErrLoadingEnvFile = "error loading environment file"

	ErrDashboardGetCountryData       = "error getting country data"
	ErrDashboardGetCurrencyData      = "error getting currency data"
	ErrDashboardGetWeatherData       = "error getting weather data"
	ErrDashboardMergingData          = "error merging data"
	ErrDashboardFilterByRegistration = "error filtering data by registration"
	ErrDashboardCountryNotFound      = "country not found"
	ErrDashboardCountryNotMatch      = "country does not match"
)
