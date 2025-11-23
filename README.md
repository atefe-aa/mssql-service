# MSSQL Edge Gateway API

This is a Go-based edge API application designed to query a MSSQL Gateway in a private network. It allows external clients to retrieve barcode records or execute raw queries on the database through a central gateway. The app is configurable via `.env` and is structured to easily add new APIs.

---

## Features

* **Settings UI** for Gateway URL and server port.
* **Test connection** to the gateway via `/settings/test-connection`.
* **Barcode API** to fetch main or child barcode records.
* **Health check endpoint** for monitoring.
* **Consistent API response format**.
* Expandable: easily add new endpoints that interact with the gateway or provide additional functionality.

---

## Requirements

* Go 1.20+
* MSSQL Gateway (running in the private network)
* `.env` file with configuration

---

## Configuration

Create a `.env` file in the project root:

```env
GATEWAY_URL=http://your-gateway-url:8080
SERVER_PORT=8080
```

* `GATEWAY_URL` – the URL of the MSSQL gateway API
* `SERVER_PORT` – port this edge API should run on

The `Settings` page allows updating and testing these values without restarting manually (except to apply changes to `.env`).

---

## Running the App

```bash
go run main.go
```

Server will run on the configured `SERVER_PORT`.

Endpoints:

* `GET /health` – returns `{ "status": "ok" }`
* `GET /api/barcode?barcode=<value>` – fetch barcode records
* `POST /settings/test-connection` – test the configured gateway URL
* `POST /settings/save` – save new configuration

---

## Response Structure

All API responses follow a **consistent JSON format**:

### Success:

```json
{
  "success": true,
  "rows": [ /* array of records */ ],
  "count": 10
}
```

### Failure:

```json
{
  "success": false,
  "error": "error_code",
  "details": "optional detailed message"
}
```

This structure is applied across barcode queries, raw query execution, and connection testing.

---

## Logging

All logs are saved in `app.log` with `log` package. This includes startup, errors, and connection test results.

---

## License

MIT License.
