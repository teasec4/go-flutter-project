# Bank API Server

A simple REST API server for managing bank accounts, built with Go.

## Quick Links

- **[QUICKSTART.md](QUICKSTART.md)** - Getting started and API examples
- **[FEATURES.md](FEATURES.md)** - Detailed feature documentation
- **Code** - All files have complete English comments

## Features

### REST API Endpoints
```bash
GET    /account?id=1           # Get account balance
POST   /account/deposit        # Deposit money
POST   /account/withdraw       # Withdraw money
```

### Request Logging
All requests are logged with execution time:
```
[13:10:08] GET /account?id=1 - 323.334µs
[13:10:08] POST /account/deposit - 122.709µs
```

### Graceful Shutdown
- Listens for SIGINT (Ctrl+C) and SIGTERM
- Finishes in-flight requests (up to 10 seconds)
- Closes connections properly
- Full context.Context support

## Getting Started

### Run the server
```bash
cd cmd/api
go run .
```

### Test
```bash
./test.sh
```

### API Examples
```bash
# Get balance
curl http://localhost:8080/account?id=1

# Deposit 500
curl -X POST http://localhost:8080/account/deposit \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":500}'

# Withdraw 200
curl -X POST http://localhost:8080/account/withdraw \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":200}'
```

## Project Structure

```
server_go/
├── cmd/api/
│   ├── main.go           # Server and graceful shutdown
│   └── types.go          # Response types
├── internal/
│   ├── account/          # Account interface and implementation
│   ├── bank/             # Thread-safe bank
│   ├── handler/          # HTTP handlers and middleware
│   ├── middleware/       # Additional middleware
│   └── tool/             # Database interface and mock
├── QUICKSTART.md         # Examples and FAQ
├── FEATURES.md           # Detailed documentation
└── test.sh               # Test script
```

## Technology Stack

- **Go 1.25.1**
- **Chi Router** - HTTP routing
- **Standard Library** - HTTP, context, signals

## Key Concepts

### Thread-safe Bank
```go
type Bank struct {
    accounts map[string]account.Account
    mu sync.Mutex  // Protects against race conditions
}
```

### Context for Shutdown
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
server.Shutdown(ctx)
```

### Middleware Pattern
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do something
        next.ServeHTTP(w, r)
        // Do something after
    })
}
```

## What to Learn

1. **main.go** - Signal handling, graceful shutdown
2. **handler.go** - HTTP handlers, middleware, request/response
3. **bank.go** - Mutex, thread-safety patterns
4. **account.go** - Interface design, validation

## API Responses

### Success (200 OK)
```json
{
  "accountId": "1",
  "balance": 1500
}
```

### Error (400, 404, etc)
```json
{
  "error": "insufficient balance"
}
```

## Configuration

### Server Port
`cmd/api/main.go` → line with `":8080"`

### Graceful Shutdown Timeout
`cmd/api/main.go` → `context.WithTimeout(..., 10*time.Second)`

### Server Timeouts
`cmd/api/main.go` → `http.Server` settings:
- ReadTimeout: 15s
- WriteTimeout: 15s
- IdleTimeout: 60s

## Shutdown Signals

Server gracefully stops on:
- **Ctrl+C** (SIGINT)
- **kill -TERM** (SIGTERM)
- **kill -INT** (SIGINT)

## Test Accounts

| ID  | Balance |
|-----|---------|
| "1" | 1000    |
| "2" | 2000    |

## Learning Goals

Use this project to understand:
- HTTP server in Go
- Routing and middleware
- Graceful shutdown with context
- Signal handling (SIGINT, SIGTERM)
- Thread-safe concurrent access (mutex)
- Interface design
- Error handling
- JSON marshaling/unmarshaling

## Useful Links

- [Go HTTP Package](https://golang.org/pkg/net/http/)
- [Chi Router](https://github.com/go-chi/chi)
- [Context Package](https://golang.org/pkg/context/)
- [Sync Package](https://golang.org/pkg/sync/)

---

Start with `QUICKSTART.md` for a quick overview, then move to `FEATURES.md` for details.
