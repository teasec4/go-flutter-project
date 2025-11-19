# Go Flutter Banking App

A simple banking application with a Flutter frontend and Go backend.

## Features

- View account information (balance, account number, owner name)
- Perform deposits and withdrawals
- Real-time account data from backend API

## Tech Stack

- **Frontend**: Flutter (iOS, macOS, Android, Web, Windows, Linux)
- **Backend**: Go HTTP server
- **Communication**: REST API

## Project Structure

- `frontend_flutter/` - Flutter application
- `server/` - Go backend server (separate repo/folder)

## Getting Started

### Frontend

```bash
cd frontend_flutter
flutter pub get
flutter run
```

### Backend

```bash
cd server
go run main.go
```

Server runs on `http://localhost:8080`

## API

- `GET /balance?id={accountId}` - Get account information

Response:
```json
{
  "accountId": "1",
  "balance": 1000.0
}
```

## Development Notes

- Frontend connects to backend via `192.168.5.10:8080` (update IP as needed)
- macOS requires network client entitlements for local network access
