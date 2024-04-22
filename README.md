# Countries Dashboard Service

This is a REST web application in Golang that provides the client with the ability to configure information dashboards
that are dynamically populated when requested. The dashboard configurations are saved in the service in a persistent
way,
and populated based on external services. It also includes a simple notification service that can listen to specific
events.
The application will be dockerized and deployed using an IaaS system.

[TOC]

## Endpoints

The service provides the following endpoints:

```plaintext
/dashboard/v1/registrations/
/dashboard/v1/dashboards/
/dashboard/v1/notifications/
/dashboard/v1/status/
```

---

### Registrations

The initial endpoint focuses on the management of dashboard configurations that can later be used via the `dashboards`
endpoint.

#### Register new dashboard configuration

Manages the registration of new dashboard configurations that indicate which information is to be shown for registered
dashboards (via the `dashboards` endpoint -- see later). This includes weather, country and currency exchange
information.

##### Request

```http
Method: POST
Path: /dashboard/v1/registrations/
Content type: application/json
```

Body example:

```json lines
{
  // Indicates country name (alternatively to ISO code, i.e., country name can be empty if ISO code field is filled and vice versa)
  "country": "Norway",
  // Indicates two-letter ISO code for country (alternatively to country name)
  "isoCode": "NO",
  "features": {
    // Indicates whether temperature in degree Celsius is shown
    "temperature": true,
    // Indicates whether precipitation (rain, showers and snow) is shown
    "precipitation": true,
    // Indicates whether the name of the capital is shown
    "capital": true,
    // Indicates whether country coordinates are shown
    "coordinates": true,
    // Indicates whether population is shown
    "population": true,
    // Indicates whether land area size is shown
    "area": true,
    // Indicates which exchange rates (to target currencies) relative to the base currency of the registered country (in this case NOK for Norway) are shown
    "targetCurrencies": [
      "EUR",
      "USD",
      "SEK"
    ]
  }
}
```

##### Response

The response to the POST request on the endpoint stores the configuration on the server and returns the associated ID.
In the example below, it is the ID `1`. Responses show be encoded in the above-mentioned JSON format, with the
`lastChange` field highlighting the last change to the configuration (including updates via `PUT` - see later)

* Content type: application/json
* Status code: Appropriate error code. Ensure to deal with errors gracefully.

Body (exemplary code for registered configuration):

```json
{
  "id": "621effa4",
  "lastChange": "2024-04-18T16:30:38.066008+02:00"
}

```

#### View a specific registered dashboard configuration

Enables retrieval of a specific registered dashboard configuration.

##### Request

The following shows a request for an individual configuration identified by its ID.

```text
Method: GET
Path: /dashboard/v1/registrations/{id}
```

* `id` is the ID associated with the specific configuration.

Example request:

```http request
/dashboard/v1/registrations/621effa4
```

##### Response

* Content type: `application/json`
* Status code: Appropriate error code. Ensure to deal with errors gracefully.

Body (exemplary code):

```json
{
  "id": "621effa4",
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": false,
    "coordinates": false,
    "population": true,
    "area": true,
    "targetCurrencies": [
      "EUR",
      "USD"
    ]
  },
  "lastChange": "2024-04-18T14:30:38.066008Z"
}
```

#### View all registered dashboard configurations

Enables retrieval of all registered dashboard configurations.

##### Request

A `GET` request to the endpoint should return all registered configurations including IDs and timestamps of last change.

```text
Method: GET
Path: /dashboard/v1/registrations/
```

##### Response

* Content type: `application/json`
* Status code: Appropriate error code. Ensure to deal with errors gracefully.

Body (exemplary code):

```json lines
[
  {
    "id": "994175d9",
    "country": "Sweden",
    "isoCode": "SE",
    "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": true,
      "targetCurrencies": [
        "EUR",
        "USD",
        "NOK"
      ]
    },
    "lastChange": "2024-04-16T11:31:05.258726Z"
  },
  {
    "id": "a60a9989",
    "country": "Sweden",
    "isoCode": "SE",
    "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": true,
      "targetCurrencies": [
        "EUR",
        "USD",
        "NOK"
      ]
    },
    "lastChange": "2024-04-16T13:01:12.452834Z"
  }
]
```

The response should return a collection of return all stored configurations.

#### View header of all registered dashboard configurations

Enables retrieval the header of all registered dashboard configurations.

##### Request

The following shows a request for the header of all registered dashboard configurations.

```text
Method: HEAD
Path: /dashboard/v1/registrations/
```

##### Response

Empty response with the header information.

* Status code: Appropriate error code.
* Body: empty

#### Replace a specific registered dashboard configuration

Enables the replacing of specific registered dashboard configurations.

##### Request

The following shows a request for an updated of individual configuration identified by its ID. This update should lead
to an update of the configuration and an update of the associated timestamp (`lastChange`).

```
Method: PUT
Path: /dashboard/v1/registrations/{id}
```

* `id` is the ID associated with the specific configuration.

Example request: ```/dashboard/v1/registrations/621effa4```

Body (exemplary code):

```json lines
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    // temperature value is to be changed
    "temperature": false,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": false,
    // targetCurrencies value is to be changed
    "targetCurrencies": [
      "EUR",
      "SEK"
    ]
  }
}
```

Note that the request neither contains ID in the body (only in the URL), and neither contains the timestamp. The
timestamp should be generated on the server side.

TODO: Implement/remove this

**Advanced Task:** Implement the `PATCH` method functionality.

##### Response

This is the response to the change request.

* Status code: Appropriate error code.
* Body: empty

#### Delete a specific registered dashboard configuration

Enabling the deletion of a specific registered dashboard configuration.

##### Request

The following shows a request for deletion of an individual configuration identified by its ID. This update should lead
to a deletion of the configuration on the server.

```text
Method: DELETE
Path: /dashboard/v1/registrations/{id}
```

* `id` is the ID associated with the specific configuration.

Example request:

```http request
/dashboard/v1/registrations/621effa4
```

##### Response

This is the response to the delete request.

* Status code: Appropriate error code.
* Body: empty

---

### Dashboards

This endpoint can be used to retrieve the populated dashboards.

#### Retrieving a specific populated dashboard

##### Request

The following shows a request for an individual dashboard identified by its ID (same as the corresponding configuration
ID).

```text
Method: GET
Path: /dashboard/v1/dashboards/{id}
```

* `id` is the ID associated with the specific configuration.

Example request:

```http request
/dashboard/v1/dashboards/621effa4
```

##### Response

* Content type: `application/json`
* Status code: Appropriate error code.

Body (exemplary code):

```json lines
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "precipitation": 0,
    "capital": "Oslo",
    "coordinates": {
      "latitude": 62,
      "longitude": 10
    },
    "population": 5379475,
    "targetCurrencies": {
      "EUR": 0.085272,
      "SEK": 0.995781
    },
    "currency": {
      "name": "Norwegian krone",
      "symbol": "kr",
      "code": "NOK"
    }
  },
  "lastRetrieval": "2024-04-18T16:37:42.469867+02:00"
}
```

---

### Notifications

As an additional feature, users can register webhooks that are triggered by the service based on specified events,
specifically if a new configuration is created, changed or deleted. Users can also register for invocation events, i.e.,
when a dashboard for a given country is invoked. Users can register multiple webhooks. The registrations should survive
a service restart (i.e., be persistently stored).

#### Registration of Webhook

##### Request

```text
Method: POST
Path: /dashboard/v1/notifications/
Content type: application/json
```

The body contains

* the URL to be triggered upon event (the service that should be invoked)
* the country for which the trigger applies (if empty, it applies to any invocation)
* Events:
    * `REGISTER` - webhook is invoked if a new configuration is registered
    * `CHANGE` - webhook is invoked if configuration is modified
    * `DELETE` - webhook is invoked if configuration is deleted
    * `INVOKE` - webhook is invoked if dashboard is retrieved (i.e., populated with values)

Body (Exemplary message based on schema):

```json lines
{
  // URL to be invoked when event occurs
  "url": "https://localhost:8080/client/",
  // Country that is registered, or empty if all countries
  "country": "NO",
  // Event on which it is invoked
  "event": "INVOKE"
}
```

##### Response

The response contains the ID for the registration that can be used to see detail information or to delete the webhook
registration. The format of the ID is not prescribed, as long it is unique. Consider best practices for determining IDs.

* Content type: `application/json`
* Status code: Appropriate status code

Body (Exemplary message based on schema):

```json
{
  "id": "OIdksUDwveiwe"
}
```

#### Deletion of Webhook

Deletes a given webhook.

##### Request

```text
Method: DELETE
Path: /dashboard/v1/notifications/{id}
```

* {id} is the ID returned during the webhook registration

##### Response

TODO: Implement and update this

Implement the response according to best practices.

#### View specific registered webhook

Shows a specific webhook registration.

##### Request

```text
Method: GET
Path: /dashboard/v1/notifications/{id}
```

* `{id}` is the ID for the webhook registration

##### Response

The response is similar to the POST request body, but further includes the ID assigned by the server upon adding the
webhook.

* Content type: `application/json`

Body (Exemplary message based on schema):

```json
{
  "id": "OIdksUDwveiwe",
  "url": "https://localhost:8080/client/",
  "country": "NO",
  "event": "INVOKE"
}

```

#### View all registered webhooks

Lists all registered webhooks.

##### Request

```text
Method: GET
Path: /dashboard/v1/notifications/
```

##### Response

The response is a collection of all registered webhooks.

* Content type: `application/json`

Body (Exemplary message based on schema):

```json lines
[
  {
    "id": "OIdksUDwveiwe",
    "url": "https://localhost:8080/client/",
    "country": "NO",
    "event": "INVOKE"
  },
  {
    "webhook_id": "DiSoisivucios",
    "url": "https://localhost:8081/anotherClient/",
    // field can also be omitted if registered for all countries
    "country": "",
    "event": "REGISTER"
  },
  ...
]
```

#### Webhook Invocation (upon trigger)

When a webhook is triggered, it should send information as follows. Where multiple webhooks are triggered, the
information should be sent separately (i.e., one notification per triggered webhook). Note that for testing purposes,
this will require you to set up another service that is able to receive the invocation. During the development, consider
using https://webhook.site/ initially.

```text
Method: POST
Path: <url specified in the corresponding webhook registration>
Content type: application/json
```

Body (Exemplary message based on schema):

```json lines
{
  "id": "OIdksUDwveiwe",
  "country": "NO",
  "event": "INVOKE",
  // time at which the event occurred
  "time": "20240223 06:23"
}
```

TODO: Implement/remove this

* **Advanced Task:** Consider supporting other event types you can think of.

---

### Status

The status interface indicates the availability of all individual services this service depends on. These can be more
services than the ones specified above (if you see any need). If you include more, you can specify additional keys with
the suffix `api`. The reporting occurs based on status codes returned by the dependent services. The status interface
further provides information about the number of registered webhooks (specification below), and the uptime of the
service.

#### Request

```text
Method: GET
Path: dashboard/v1/status/
```

#### Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise.

Body:

```json lines
{
  "countries_api": "http status code for *REST Countries API*",
  "meteo_api": "http status code for *Meteo API*",
  "currency_api": "http status code for *Currency API*",
  "notification_db": "http status code for *Notification database*",
  "...": "...",
  "webhooks": "number of registered webhooks",
  "version": "v1",
  "uptime": "time in seconds from the last service restart"
}
```

TODO: Implement/remove this

Note: `<some value>` indicates placeholders for values to be populated by the service as described for the corresponding
values. Feel free to extend the output with information you deem useful to assess the status of your service.

---

## Configuration

The service can be configured using the following environment variables:

```dotenv
PORT=
TEST_PORT=
TYPE=
PROJECTID=
PRIVATEKEYID=
PRIVATEKEY=
CLIENTEMAIL=
CLIENTID=
AUTHURI=
TOKENURI=
AUTHPROVIDERX509CERTURL=
CLIENTX509CERTURL=
UNIVERSEDOMAIN=
```

See the empty .env file for an example. Most of the variables are used for Firebase authentication.

## Deployment

The service can be deployed using the following command:

```bash
docker compose up -d --build
```

### Logs

```bash
docker ps -a
```


```bash
docker logs -f <container_name>
```

## Testing

Run the run_tests.sh script to run all tests. Add +x permission to the script if needed.
Firebase installation and authentication is required for most tests, so tests will fail if not set up.
See coverage report for information about test coverage.

```bash
./run_tests.sh
```

This command was used to generate the coverage report:

```bash
clear; go clean -testcache; ./run_tests.sh; go tool cover -html=coverage.out -o coverage.html; rm coverage.out
```

External services are mocked in the tests, so no external services are required to run the tests.

### Coverage

Since firestore emulator is required for testing, one cannot run the tests without setting up firebase. Therefore, the
coverage report is not available in the repository. The coverage report can be viewed by running the following command:

```bash
open coverage.html
```

## Known issues

The tests fails unless firebase is set up. This is because the tests require the firestore emulator to be running.

## Future work

Input validation is not thoroughly implemented in the service. This should be implemented to ensure that the service is robust and
secure.

## Contact

For any questions, please contact the authors at email:

* [Erik Bjørnsen](mailto:erbj@stud.ntnu.no) `erbj@stud.ntnu.no`
* [Simon Houmb](mailto:simonhou@ntnu.no) `simonhou@ntnu.no`