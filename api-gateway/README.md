# API Gateway

API Gateway для агрегации всех микросервисов проекта CG-2025-Crutch.

## Описание

API Gateway - это единая точка входа для всех микросервисов системы. Он обрабатывает:
- User Service - управление пользователями
- Funds Service - управление финансами
- Notification Service - управление уведомлениями

Все эндпоинты, кроме регистрации и логина, требуют авторизации через JWT токен.

## Архитектура

```
api-gateway/
├── cmd/
│   └── main.go                 # Точка входа
├── internal/
│   ├── clients/                # gRPC клиенты для микросервисов
│   │   └── clients.go
│   ├── config/                 # Конфигурация
│   │   └── config.go
│   ├── handlers/               # HTTP handlers
│   │   ├── user/               # User handlers
│   │   ├── funds/              # Funds handlers
│   │   └── notifications/      # Notifications handlers
│   ├── middleware/             # Middleware
│   │   └── auth.go             # JWT авторизация
│   ├── grpc/gen/               # Сгенерированные protobuf файлы
│   └── run/                    # Логика запуска приложения
│       └── app.go
├── Dockerfile
└── go.mod
```

## Конфигурация

Создайте файл `.env` на основе `.env.example`:

```env
API_GATEWAY_HOST=
API_GATEWAY_PORT=8080

USER_SERVICE_HOST=user-service
USER_SERVICE_PORT=50051
USER_SERVICE_TIMEOUT=10s

FUNDS_SERVICE_HOST=funds-service
FUNDS_SERVICE_PORT=50052
FUNDS_SERVICE_TIMEOUT=10s

NOTIFICATION_SERVICE_HOST=notification-service
NOTIFICATION_SERVICE_PORT=50053
NOTIFICATION_SERVICE_TIMEOUT=10s
```

## API Endpoints

### Health Check
```
GET /health
```
Проверка работоспособности сервиса.

### User Service

#### Public Endpoints (не требуют авторизации)

##### Регистрация пользователя
```
POST /api/v1/users/register
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securepassword123",
  "first_name": "John",
  "second_name": "Doe",
  "age": 25,
  "salary": 50000.0,
  "work_sphere_id": 1
}
```

##### Логин пользователя
```
POST /api/v1/users/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securepassword123"
}

Response:
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_at": 1234567890,
  "user": { ... }
}
```

#### Secured Endpoints (требуют Authorization header)

##### Обновить пользователя
```
PUT /api/v1/users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "john_doe",
  "first_name": "John",
  "second_name": "Doe",
  "age": 26,
  "salary": 55000.0,
  "work_sphere_id": 1
}
```

##### Получить пользователя по ID
```
GET /api/v1/users/:id
Authorization: Bearer <token>
```

##### Получить пользователя по username
```
GET /api/v1/users/username/:username
Authorization: Bearer <token>
```

### Funds Service (все требуют авторизации)

#### Transactions

##### Создать транзакцию
```
POST /api/v1/funds/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": 1,
  "type": "expense",
  "amount": 100.50,
  "title": "Grocery shopping",
  "description": "Weekly groceries",
  "transaction_date": "2024-12-06"
}
```

##### Получить транзакцию по ID
```
GET /api/v1/funds/transactions/:id
Authorization: Bearer <token>
```

##### Получить все транзакции пользователя
```
GET /api/v1/funds/transactions?limit=10&offset=0
Authorization: Bearer <token>
```

##### Получить транзакции пользователя за период
```
GET /api/v1/funds/transactions/period?days=30&limit=10&offset=0
Authorization: Bearer <token>
```

##### Обновить транзакцию
```
PUT /api/v1/funds/transactions/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": 1,
  "type": "expense",
  "amount": 120.50,
  "title": "Grocery shopping",
  "description": "Weekly groceries + snacks",
  "transaction_date": "2024-12-06"
}
```

##### Удалить транзакцию
```
DELETE /api/v1/funds/transactions/:id
Authorization: Bearer <token>
```

#### Categories

##### Получить все категории
```
GET /api/v1/funds/categories
Authorization: Bearer <token>
```

##### Получить категории по типу
```
GET /api/v1/funds/categories/type/:type
Authorization: Bearer <token>

:type - "income" или "expense"
```

##### Получить категорию по ID
```
GET /api/v1/funds/categories/:id
Authorization: Bearer <token>
```

#### Balance

##### Получить баланс пользователя
```
GET /api/v1/funds/balance
Authorization: Bearer <token>
```

### Notification Service (все требуют авторизации)

##### Получить VAPID ключ
```
GET /api/v1/notifications/vapid-key
Authorization: Bearer <token>
```

##### Подписаться на уведомления
```
POST /api/v1/notifications/subscribe
Authorization: Bearer <token>
Content-Type: application/json

{
  "endpoint": "https://fcm.googleapis.com/fcm/send/...",
  "p256dh": "BKy...",
  "auth": "YM5..."
}
```

## Авторизация

Все защищенные эндпоинты требуют JWT токен в заголовке:

```
Authorization: Bearer <access_token>
```

Токен получается при логине через `/api/v1/users/login`. При каждом запросе к защищенному эндпоинту:
1. API Gateway извлекает токен из заголовка Authorization
2. Валидирует токен через User Service (метод ValidateToken)
3. Извлекает user_id из токена
4. Добавляет user_id в контекст запроса для дальнейшего использования

## Запуск

### Локально
```bash
go run cmd/main.go
```

### Docker
```bash
docker build -t api-gateway .
docker run -p 8080:8080 --env-file .env api-gateway
```

### Docker Compose
API Gateway интегрирован в общий `docker-compose.yml` проекта.

## Middleware

### Auth Middleware
- Проверяет наличие Authorization header
- Валидирует формат токена (Bearer <token>)
- Вызывает User Service для валидации токена
- Добавляет user_id в контекст запроса

### CORS
- Разрешены все origins (*)
- Разрешенные методы: GET, POST, PUT, DELETE, OPTIONS
- Разрешенные заголовки: Origin, Content-Type, Accept, Authorization

### Logger
- Логирование всех HTTP запросов

### Recover
- Обработка паник для graceful recovery

## Обработка ошибок

API Gateway возвращает ошибки в формате:
```json
{
  "error": "error message"
}
```

Коды ошибок:
- 400 - Bad Request (невалидные данные)
- 401 - Unauthorized (отсутствует или невалидный токен)
- 403 - Forbidden (нет прав доступа)
- 404 - Not Found (ресурс не найден)
- 500 - Internal Server Error (внутренняя ошибка)

## Технологии

- Go 1.24.3
- Fiber v2 - HTTP framework
- gRPC - коммуникация с микросервисами
- Zap - структурированное логирование

