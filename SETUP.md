# Crypto Bank Project - Setup Guide

Полное руководство по установке и запуску банковского приложения с поддержкой криптовалюты.

## Требования

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.23+ (для локальной разработки)
- Make (опционально, для удобства)

## Быстрый старт

### 1. Клонирование и подготовка

```bash
cd crypto-bank-project
cp .env.example .env
```

### 2. Запуск всех сервисов

```bash
# С использованием make
make build
make up

# Или напрямую через docker-compose
docker-compose build
docker-compose up -d
```

### 3. Проверка статуса

```bash
# Проверка всех сервисов
make ps

# Или
docker-compose ps
```

### 4. Проверка health endpoints

```bash
# Проверка всех сервисов
make api-health

# Или по отдельности
curl http://localhost:8080/health  # Bank Service
curl http://localhost:8085/health  # Exchange Service
curl http://localhost:8082/health  # Analytics Service
curl http://localhost:8083/health  # Notification Service
```

## Структура проекта

```
crypto-bank-project/
├── bank-service/           # Основной банковский сервис
│   ├── cmd/server/        # Точка входа
│   ├── internal/          # Внутренняя логика
│   │   ├── models/       # Модели данных
│   │   ├── repositories/ # Работа с БД
│   │   ├── services/     # Бизнес-логика
│   │   ├── handlers/     # HTTP handlers
│   │   └── middleware/   # Middleware
│   ├── migrations/       # Миграции БД
│   └── pkg/             # Утилиты
├── exchange-service/      # Сервис обмена валют (gRPC)
├── analytics-service/     # Сервис аналитики
├── notification-service/  # Сервис уведомлений
├── docker-compose.yml    # Конфигурация Docker
└── Makefile             # Команды для управления
```

## Сервисы и порты

| Сервис | Порт | Описание |
|--------|------|----------|
| Bank Service | 8080 | Основное REST API |
| Exchange Service (HTTP) | 8085 | HTTP API для курсов |
| Exchange Service (gRPC) | 9090 | gRPC API для курсов |
| Analytics Service | 8082 | API аналитики |
| Notification Service | 8083 | API уведомлений |
| PostgreSQL | 5432 | База данных |
| RabbitMQ | 5672 | Очередь сообщений |
| RabbitMQ Management | 15672 | UI управления (guest/guest) |
| Prometheus | 9091 | Метрики |
| Grafana | 3000 | Дашборды (admin/admin) |
| Zipkin | 9411 | Трейсинг |

## API Endpoints

### Bank Service (http://localhost:8080)

#### Users
- `POST /api/v1/users` - Создать пользователя
- `GET /api/v1/users` - Получить всех пользователей
- `GET /api/v1/users/:id` - Получить пользователя
- `PUT /api/v1/users/:id` - Обновить пользователя
- `DELETE /api/v1/users/:id` - Удалить пользователя

#### Accounts
- `POST /api/v1/accounts` - Создать фиатный счет
- `GET /api/v1/accounts/:id` - Получить счет
- `GET /api/v1/users/:user_id/accounts` - Получить все счета пользователя
- `GET /api/v1/accounts/:id/balance` - Получить баланс

#### Crypto Wallets
- `POST /api/v1/wallets` - Создать криптокошелек
- `GET /api/v1/wallets/:id` - Получить кошелек
- `GET /api/v1/users/:user_id/wallets` - Получить все кошельки пользователя
- `GET /api/v1/wallets/:id/balance` - Получить баланс

#### Transactions
- `POST /api/v1/transactions/transfer` - Перевод между счетами
- `POST /api/v1/transactions/deposit` - Пополнение счета
- `POST /api/v1/transactions/withdraw` - Снятие со счета
- `GET /api/v1/transactions/:id` - Получить транзакцию
- `GET /api/v1/users/:user_id/transactions` - История транзакций

#### Exchanges
- `POST /api/v1/exchanges/crypto-to-fiat` - Обменять крипту на фиат
- `POST /api/v1/exchanges/fiat-to-crypto` - Обменять фиат на крипту
- `GET /api/v1/exchanges/:id` - Получить обмен
- `GET /api/v1/users/:user_id/exchanges` - История обменов

### Analytics Service (http://localhost:8082)

- `GET /api/v1/statistics` - Получить статистику

### Notification Service (http://localhost:8083)

- `GET /api/v1/notifications` - Получить все уведомления
- `GET /api/v1/notifications/:user_id` - Получить уведомления пользователя

## Примеры использования

### 1. Создание пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

### 2. Создание фиатного счета

```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_ID_HERE",
    "currency": "USD"
  }'
```

### 3. Создание криптокошелька

```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_ID_HERE",
    "crypto_type": "BTC"
  }'
```

### 4. Пополнение счета

```bash
curl -X POST http://localhost:8080/api/v1/transactions/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "ACCOUNT_ID_HERE",
    "amount": 1000.00
  }'
```

### 5. Обмен криптовалюты на фиат

```bash
curl -X POST http://localhost:8080/api/v1/exchanges/crypto-to-fiat \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_ID_HERE",
    "from_wallet_id": "WALLET_ID_HERE",
    "to_account_id": "ACCOUNT_ID_HERE",
    "crypto_amount": 0.01
  }'
```

## Мониторинг

### Prometheus
```bash
make open-prometheus
# или
open http://localhost:9091
```

### Grafana
```bash
make open-grafana
# или
open http://localhost:3000
# Логин: admin / admin
```

### Zipkin (трейсинг)
```bash
make open-zipkin
# или
open http://localhost:9411
```

### RabbitMQ Management
```bash
make open-rabbitmq
# или
open http://localhost:15672
# Логин: guest / guest
```

## Логи

```bash
# Все сервисы
make logs

# Конкретный сервис
make logs-bank
make logs-exchange
make logs-analytics
make logs-notification

# Или через docker-compose
docker-compose logs -f bank-service
```

## Остановка и очистка

```bash
# Остановить все сервисы
make down

# Полная очистка (удаление volumes и images)
make clean
```

## Разработка

### Локальный запуск (без Docker)

1. Запустите инфраструктуру:
```bash
docker-compose up -d postgres rabbitmq zipkin prometheus grafana
```

2. Запустите сервисы локально:
```bash
# В разных терминалах
make dev-exchange
make dev-analytics
make dev-notification
make dev-bank
```

### Запуск тестов

```bash
make test
```

## Технологии

- **Fiber** - веб-фреймворк для Go
- **Squirrel** - SQL query builder
- **Uber Zap** - структурированное логирование
- **Goose** - миграции базы данных
- **PostgreSQL** - реляционная БД
- **RabbitMQ** - брокер сообщений
- **gRPC** - межсервисное взаимодействие
- **Prometheus** - сбор метрик
- **Grafana** - визуализация метрик
- **Zipkin** - распределенная трассировка

## Troubleshooting

### Проблемы с подключением к БД

```bash
# Проверьте статус PostgreSQL
docker-compose ps postgres

# Посмотрите логи
docker-compose logs postgres
```

### Проблемы с RabbitMQ

```bash
# Проверьте статус
docker-compose ps rabbitmq

# Посмотрите логи
docker-compose logs rabbitmq

# Откройте Management UI
open http://localhost:15672
```

### Пересоздание контейнеров

```bash
docker-compose down
docker-compose up -d --force-recreate
```

### Очистка volumes

```bash
docker-compose down -v
docker-compose up -d
```

## Поддержка

Для вопросов и предложений создавайте issue в репозитории.

