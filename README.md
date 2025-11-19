# MyApp - Go API with MSSQL

A clean, simple Go API for reading data from MSSQL database with support for both Windows and SQL authentication.

## Project Structure

```
myapp/
├── cmd/api/main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── database/             # Database connection
│   ├── handlers/             # HTTP handlers
│   ├── models/               # Data models
│   └── repository/           # Database queries
├── .env.example              # Environment variables template
├── go.mod
└── README.md
```

## Setup

1. **Initialize the project:**
   ```bash
   go mod download
   ```

2. **Configure environment variables:**
   
   Copy `.env.example` to `.env` and configure:
   
   **For Windows Authentication:**
   ```env
   DB_SERVER=localhost
   DB_PORT=1433
   DB_NAME=your_database
   DB_AUTH_MODE=windows
   SERVER_PORT=8080
   ```
   
   **For SQL Authentication:**
   ```env
   DB_SERVER=localhost
   DB_PORT=1433
   DB_NAME=your_database
   DB_AUTH_MODE=sql
   DB_USER=your_username
   DB_PASSWORD=your_password
   SERVER_PORT=8080
   ```

3. **Run the application:**
   ```bash
   # On Windows with environment variables
   go run cmd/api/main.go
   
   # Or set environment variables inline (Linux/Mac)
   DB_NAME=mydb DB_AUTH_MODE=windows go run cmd/api/main.go
   ```

## API Endpoints

- `GET /health` - Health check
- `GET /api/users` - Get all users
- `GET /api/users/by-id?id=1` - Get user by ID
- `GET /api/products` - Get all products

## Adding New Features

### Adding a New API Endpoint

1. **Create model** in `internal/models/models.go`
2. **Add repository method** in `internal/repository/repository.go`
3. **Create handler** in `internal/handlers/handlers.go`
4. **Register route** in `cmd/api/main.go`

### Example: Adding a "Categories" endpoint

**Step 1: Add model**
```go
// internal/models/models.go
type Category struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
```

**Step 2: Add repository method**
```go
// internal/repository/repository.go
func (r *Repository) GetCategories(ctx context.Context) ([]models.Category, error) {
    query := `SELECT id, name FROM categories`
    // ... implementation
}
```

**Step 3: Add handler**
```go
// internal/handlers/handlers.go
func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := h.repo.GetCategories(r.Context())
    // ... implementation
}
```

**Step 4: Register route**
```go
// cmd/api/main.go
mux.HandleFunc("/api/categories", handler.GetCategories)
```

## Configuration Options

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| DB_SERVER | Database server address | localhost | No |
| DB_PORT | Database port | 1433 | No |
| DB_NAME | Database name | - | Yes |
| DB_AUTH_MODE | Authentication mode (`windows` or `sql`) | windows | No |
| DB_USER | Username for SQL auth | - | If using SQL auth |
| DB_PASSWORD | Password for SQL auth | - | If using SQL auth |
| SERVER_PORT | HTTP server port | 8080 | No |

## Building for Production

```bash
# Build executable
go build -o myapp.exe cmd/api/main.go

# Run with environment variables
set DB_NAME=production_db
set DB_AUTH_MODE=windows
set SERVER_PORT=9000
myapp.exe
```

## Notes

- All database operations are read-only
- Connection pooling is configured automatically (max 25 connections)
- Uses context for request timeouts and cancellation
- Structured logging for debugging