# Bank API Server - Quick Start

## 🚀 Running the Server

```bash
cd cmd/api
go run .
```

Expected output:
```
2025/11/21 13:09:49 API server started on port :8080
[13:09:50] GET /account?id=1
[13:09:50] GET /account?id=1 - 323.334µs
```

## 📡 API Endpoints

### 1. Get Account Balance
```bash
GET /account?id=1

curl -X GET "http://localhost:8080/account?id=1"
```

**Response (200 OK):**
```json
{
  "accountId": "1",
  "balance": 1000
}
```

### 2. Deposit Funds
```bash
POST /account/deposit

curl -X POST "http://localhost:8080/account/deposit" \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":500}'
```

**Response (200 OK):**
```json
{
  "accountId": "1",
  "balance": 1500
}
```

### 3. Withdraw Funds
```bash
POST /account/withdraw

curl -X POST "http://localhost:8080/account/withdraw" \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":200}'
```

**Response (200 OK):**
```json
{
  "accountId": "1",
  "balance": 1300
}
```

## 🔍 Default Accounts

The server initializes with two test accounts:

| Account ID | Initial Balance |
|-----------|-----------------|
| "1"       | 1000            |
| "2"       | 2000            |

## 📋 Error Responses

### Invalid Account ID
```bash
curl -X GET "http://localhost:8080/account?id=999"
```
**Response (404 Not Found):**
```json
{
  "error": "account not found"
}
```

### Missing ID Parameter
```bash
curl -X GET "http://localhost:8080/account"
```
**Response (400 Bad Request):**
```json
{
  "error": "id parameter required"
}
```

### Insufficient Balance
```bash
curl -X POST "http://localhost:8080/account/withdraw" \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":10000}'
```
**Response (400 Bad Request):**
```json
{
  "error": "insufficient balance"
}
```

### Invalid Amount
```bash
curl -X POST "http://localhost:8080/account/deposit" \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":-100}'
```
**Response (400 Bad Request):**
```json
{
  "error": "deposit amount must be greater than 0"
}
```

## 🛑 Graceful Shutdown

### Method 1: Ctrl+C
```bash
./api
^C  # Press Ctrl+C
# Server gracefully shuts down
```

### Method 2: Kill Signal
```bash
# Get server PID
ps aux | grep api

# Send SIGTERM
kill -TERM <pid>
```

**Shutdown output:**
```
Received signal: interrupt, starting graceful shutdown...
Server gracefully shut down
```

## 📊 Logging Middleware

Every request to `/account` routes is logged with timestamp and duration:

```
[13:10:08] GET /account?id=1
[13:10:08] GET /account?id=1 - 323.334µs
[13:10:08] POST /account/deposit
[13:10:08] POST /account/deposit - 122.709µs
```

**Timing breakdown:**
- `µs` = microseconds (< 1ms)
- `ms` = milliseconds (1000µs)
- `s` = seconds (1000ms)

## 🧪 Testing Script

```bash
#!/bin/bash

# Start server in background
./api &
SERVER_PID=$!
sleep 1

# Test 1: Get balance
echo "Test 1: Get balance"
curl -X GET "http://localhost:8080/account?id=1" | jq .

# Test 2: Deposit
echo "Test 2: Deposit 500"
curl -X POST "http://localhost:8080/account/deposit" \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":500}' | jq .

# Test 3: Get new balance
echo "Test 3: Get new balance"
curl -X GET "http://localhost:8080/account?id=1" | jq .

# Test 4: Withdraw
echo "Test 4: Withdraw 200"
curl -X POST "http://localhost:8080/account/withdraw" \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":200}' | jq .

# Shutdown
kill $SERVER_PID
```

## 📂 Project Structure

```
server_go/
├── cmd/api/
│   ├── main.go           # Server entry point + graceful shutdown
│   └── types.go          # Response types
├── internal/
│   ├── account/
│   │   └── account.go    # Account interface & implementation
│   ├── bank/
│   │   └── bank.go       # Bank & thread-safe account management
│   ├── handler/
│   │   └── handler.go    # HTTP handlers + logging middleware
│   ├── middleware/
│   │   └── middleware.go # Additional middleware (JSONContent, etc.)
│   ├── tool/
│   │   ├── database.go   # Database interface
│   │   └── mockdb.go     # Mock database implementation
│   └── user/
│       └── user.go       # User-related types
├── go.mod
└── go.sum
```

## 💡 Key Files to Understand

1. **`cmd/api/main.go`** - Start here to understand:
   - Server configuration
   - Graceful shutdown with context
   - Signal handling

2. **`internal/handler/handler.go`** - Learn about:
   - HTTP handlers
   - Request/response types
   - Logging middleware

3. **`internal/bank/bank.go`** - Understanding:
   - Thread-safe account management
   - Mutex usage
   - Concurrency patterns

4. **`internal/account/account.go`** - Interface design:
   - Deposit/withdraw operations
   - Validation logic
   - Error handling

## 🔗 Related Documentation

- See `FEATURES.md` for detailed feature documentation
- Go to `internal/` for fully commented code
- Check each file header for package documentation

## ❓ FAQ

**Q: How do I change the server port?**
A: Edit `cmd/api/main.go`, change `":8080"` to your desired port

**Q: How long does graceful shutdown wait?**
A: 10 seconds by default, configurable in `main()`

**Q: What happens if a request takes longer than 10s?**
A: The server will forcefully close the connection after 10s

**Q: How do I add more accounts?**
A: Edit `internal/bank/bank.go` in the `New()` function

**Q: Can I use a real database?**
A: Yes, implement the `DatabaseInterface` in `internal/tool/database.go`
