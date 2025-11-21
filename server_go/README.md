# Bank API Server

Простой REST API сервер для управления банковскими счетами, написанный на Go.

## 📚 Документация

- **[QUICKSTART.md](QUICKSTART.md)** - Быстрый старт и примеры API
- **[FEATURES.md](FEATURES.md)** - Подробное описание функций
- **Код** - Каждый файл имеет полные англ комментарии

## ✨ Основные фичи

### 1️⃣ REST API для работы со счетами
```bash
GET    /account?id=1           # Получить баланс
POST   /account/deposit        # Положить деньги
POST   /account/withdraw       # Снять деньги
```

### 2️⃣ Логирование middleware
Все запросы логируются с временем выполнения:
```
[13:10:08] GET /account?id=1 - 323.334µs
[13:10:08] POST /account/deposit - 122.709µs
```

### 3️⃣ Graceful shutdown
- Слушает сигналы SIGINT (Ctrl+C) и SIGTERM
- Завершает in-flight запросы (до 10 секунд)
- Закрывает соединения чисто
- Полная поддержка `context.Context`

## 🚀 Быстрый старт

### Запуск сервера
```bash
cd cmd/api
go run .
```

### Тестирование
```bash
./test.sh
```

### Примеры API
```bash
# Получить баланс
curl http://localhost:8080/account?id=1

# Положить 500
curl -X POST http://localhost:8080/account/deposit \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":500}'

# Снять 200
curl -X POST http://localhost:8080/account/withdraw \
  -H "Content-Type: application/json" \
  -d '{"accountId":"1","amount":200}'
```

## 📂 Структура проекта

```
server_go/
├── cmd/api/
│   ├── main.go           # Server + graceful shutdown
│   └── types.go          # Response types
├── internal/
│   ├── account/          # Account interface & impl
│   ├── bank/             # Thread-safe bank
│   ├── handler/          # HTTP handlers + middleware
│   ├── middleware/       # Additional middleware
│   └── tool/             # Database interface & mock
├── QUICKSTART.md         # Примеры и FAQ
├── FEATURES.md           # Подробная документация
└── test.sh               # Тестовый скрипт
```

## 🔧 Технологии

- **Go 1.25.1**
- **Chi Router** - для маршрутизации
- **Стандартная библиотека** - HTTP, context, signals

## 💡 Ключевые концепции

### Thread-safe Bank
```go
type Bank struct {
    accounts map[string]account.Account
    mu sync.Mutex  // Защита от race conditions
}
```

### Context для shutdown
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
server.Shutdown(ctx)
```

### Middleware pattern
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do something
        next.ServeHTTP(w, r)
        // Do something after
    })
}
```

## 🧪 Что изучить в коде

1. **main.go** - Signal handling, graceful shutdown
2. **handler.go** - HTTP handlers, middleware, request/response
3. **bank.go** - Mutex, thread-safety patterns
4. **account.go** - Interface design, validation

## 📊 API Responses

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

## ⚙️ Конфигурация

### Порт сервера
`cmd/api/main.go` → строка с `":8080"`

### Timeout graceful shutdown
`cmd/api/main.go` → `context.WithTimeout(..., 10*time.Second)`

### Timeouts для сервера
`cmd/api/main.go` → `http.Server` структура:
- ReadTimeout: 15s
- WriteTimeout: 15s
- IdleTimeout: 60s

## 🚫 Shutdown сигналы

Сервер корректно завершается по:
- **Ctrl+C** (SIGINT)
- **kill -TERM** (SIGTERM)
- **kill -INT** (SIGINT)

## 📝 Тестовые аккаунты

| ID  | Баланс |
|-----|--------|
| "1" | 1000   |
| "2" | 2000   |

## 🎯 Практика

Используй этот проект чтобы понять:
- ✅ HTTP server в Go
- ✅ Router и middleware
- ✅ Graceful shutdown с context
- ✅ Signal handling (SIGINT, SIGTERM)
- ✅ Thread-safe concurrent access (mutex)
- ✅ Interface design patterns
- ✅ Error handling
- ✅ JSON marshaling/unmarshaling

## 🔗 Полезные ссылки

- [Go HTTP Package](https://golang.org/pkg/net/http/)
- [Chi Router](https://github.com/go-chi/chi)
- [Context Package](https://golang.org/pkg/context/)
- [Sync Package](https://golang.org/pkg/sync/)

## 📝 Обновления

Все файлы содержат полные английские комментарии для обучения.

Начни с `QUICKSTART.md` для быстрого старта, затем переходи к `FEATURES.md` для деталей.

---

**Учись на этом коде!** Каждый комментарий есть для понимания. 🎓
