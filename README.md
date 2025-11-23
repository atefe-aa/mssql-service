# MSSQL Edge Gateway API

This is a Go-based edge API application designed to communicate with an MSSQL Gateway inside a private network. It exposes controlled and secure endpoints for retrieving barcode records or running raw queries through the gateway. The application loads configuration from `.env`, provides a Settings UI, and enforces access protection via an auto-generated API key.

---

## Features

* **Auto-generated API key** stored in `.env` on first run.
* **API key enforcement** on all `/api/*` endpoints (header: `X-API-Key`).
* **Local-only access** to the Settings UI.
* **Settings UI** for configuring the Gateway URL and server port.
* **Connection tester** via `/settings/test-connection`.
* **Barcode APIs** to fetch both main and child barcode records.
* **Health check endpoint**.
* **Unified API response format**.
* Easy expansion for additional features or gateway queries.

---

## Requirements

* Go 1.20+
* MSSQL Gateway running inside the private network
* `.env` file (auto-created if missing)

---

## Configuration

Your `.env` file may look like this:

```env
GATEWAY_URL=http://your-gateway-url:8080
SERVER_PORT=8080
API_KEY=auto-generated-on-first-run
```

### Configuration values:

* `GATEWAY_URL` – URL of the MSSQL Gateway
* `SERVER_PORT` – port this Edge API runs on
* `API_KEY` – generated automatically on first launch if not present

If the file does not contain an API key, one will be generated and appended.
The API key is required for all API routes.

---

## Security

### API Key Protection

All routes under `/api/*` require a valid API key to be supplied in:

```
X-API-Key: <your-key>
```

Requests without a valid key receive a `401 Unauthorized` response.

### Localhost-Only Settings Page

The following routes are restricted to `127.0.0.1` and `::1`:

* `/settings`
* `/settings/save`
* `/settings/test-connection`

This prevents remote configuration access.

---

## Running the App

```bash
go run main.go
```

The server will start on the port defined in `.env`.

Available endpoints:

* `GET /health` – returns `{ "status": "ok" }`
* `GET /api/barcode?barcode=<value>` – fetch barcode records (requires API key)
* `GET /api/batch/barcode?barcode=<value>` – fetch a batch of barcode records (requires API key)
* `POST /settings/test-connection` – test Gateway connectivity (Only accessable from the local machine)
* `POST /settings/save` – save updated configuration (Only accessable from the local machine)

---

## API Response Structure

All API responses follow a unified JSON format.

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

---

## Logging

All logs are written to `app.log`, including startup events, API errors, connection tests, and configuration issues.

---

## License

MIT License.
