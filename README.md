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
  "id": 1,
  "lastChange": "2024-02-29 12:31"
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
/dashboard/v1/registrations/1
```

##### Response

* Content type: `application/json`
* Status code: Appropriate error code. Ensure to deal with errors gracefully.

Body (exemplary code):

```json
{
  "id": 1,
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": false,
    "targetCurrencies": [
      "EUR",
      "USD",
      "SEK"
    ]
  },
  "lastChange": "20240229 14:07"
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
    "id": 1,
    "country": "Norway",
    "isoCode": "NO",
    "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": false,
      "targetCurrencies": [
        "EUR",
        "USD",
        "SEK"
      ]
    },
    "lastChange": "20240229 14:07"
  },
  {
    "id": 2,
    "country": "Denmark",
    "isoCode": "DK",
    "features": {
      "temperature": false,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": false,
      "area": true,
      "targetCurrencies": [
        "NOK",
        "MYR",
        "JPY",
        "EUR"
      ]
    },
    "lastChange": "20240224 08:27"
  },
  ...
]
```

The response should return a collection of return all stored configurations.

TODO: Implement/remove this

**Advanced Task:** Implement the `HEAD` method functionality (only return the header, not the body).

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

Example request: ```/dashboard/v1/registrations/1```

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
/dashboard/v1/registrations/1
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
/dashboard/v1/dashboards/1
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
    // Mean temperature across all forecasted temperature values for country's coordinates
    "temperature": -1.2,
    // Mean precipitation across all returned precipitation values
    "precipitation": 0.80,
    // Capital: Where multiple values exist, take the first
    "capital": "Oslo",
    // Those are the country geocoordinates
    "coordinates": {
      "latitude": 62.0,
      "longitude": 10.0
    },
    "population": 5379475,
    "area": 323802.0,
    // this is the current NOK to EUR exchange rate (where multiple currencies exist for a given country, take the first)
    "targetCurrencies": {
      "EUR": 0.087701435,
      "USD": 0.095184741,
      "SEK": 0.97827275
    }
  },
  // this should be the current time (i.e., the time of retrieval)
  "lastRetrieval": "20240229 18:15"
}
```

Note: While it would, in principle, be easy to request everything for all registered dashboard configurations
(i.e., `GET` request on `/dashboards/` endpoint), be mindful of the services we are using. We thus only allow for one
dashboard retrieval.

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

## Testing

Run the following command to run the tests:

```bash
go test ./...
```

## Known issues

## Future work

## Contact

For any questions, please contact the authors at email:

* [Erik Bjørnsen](mailto:erbj@stud.ntnu.no) `erbj@stud.ntnu.no`
* [Simon Houmb](mailto:simonhou@ntnu.no) `simonhou@ntnu.no`