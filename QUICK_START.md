# 🚀 Быстрый старт - Crypto Bank Project

## Запуск за 3 команды

```bash
# 1. Перейдите в директорию проекта
cd crypto-bank-project

# 2. Соберите все сервисы
docker-compose build

# 3. Запустите все сервисы
docker-compose up -d
```

## Проверка работы

Откройте в браузере:
- **Bank API**: http://localhost:8080/health
- **Analytics**: http://localhost:8082/api/v1/statistics
- **Grafana**: http://localhost:3000 (admin/admin)
- **RabbitMQ**: http://localhost:15672 (guest/guest)

## Тестовые данные

В базе данных уже есть тестовые пользователи и счета после выполнения миграций.

## Пример использования API

### 1. Получить всех пользователей
```bash
curl http://localhost:8080/api/v1/users | jq
```

### 2. Создать новый фиатный счет
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "currency": "USD"
  }' | jq
```

### 3. Создать криптокошелек
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "crypto_type": "ETH"
  }' | jq
```

### 4. Пополнить счет
```bash
curl -X POST http://localhost:8080/api/v1/transactions/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "YOUR_ACCOUNT_ID",
    "amount": 5000.00
  }' | jq
```

### 5. Обменять USD на BTC
```bash
curl -X POST http://localhost:8080/api/v1/exchanges/fiat-to-crypto \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "from_account_id": "YOUR_USD_ACCOUNT_ID",
    "to_wallet_id": "YOUR_BTC_WALLET_ID",
    "fiat_amount": 100.00
  }' | jq
```

## Полная документация

Смотрите [SETUP.md](SETUP.md) для подробной документации.

## Остановка

```bash
docker-compose down
```

## Архитектура

```
┌─────────────────────────────────────────────────┐
│           Crypto Bank Application               │
├─────────────────────────────────────────────────┤
│                                                 │
│  ┌──────────────┐         ┌──────────────┐    │
│  │ Bank Service │◄───────►│  PostgreSQL  │    │
│  │  (REST API)  │         │   Database   │    │
│  └──────┬───────┘         └──────────────┘    │
│         │                                       │
│         │                  ┌──────────────┐    │
│         ├─────────────────►│   RabbitMQ   │    │
│         │                  │    Events    │    │
│         │                  └──────┬───────┘    │
│         │                         │             │
│         │         ┌───────────────┼────────┐   │
│         │         │               │        │   │
│  ┌──────▼──────┐ ┌▼──────────┐  ┌▼─────────┐ │
│  │  Exchange   │ │ Analytics │  │Notification││
│  │   Service   │ │  Service  │  │  Service   ││
│  │   (gRPC)    │ │           │  │            ││
│  └─────────────┘ └───────────┘  └────────────┘│
│                                                 │
│  ┌──────────────┐ ┌──────────────┐            │
│  │  Prometheus  │ │   Grafana    │            │
│  │   Metrics    │ │  Dashboard   │            │
│  └──────────────┘ └──────────────┘            │
│                                                 │
│  ┌──────────────┐                              │
│  │    Zipkin    │                              │
│  │   Tracing    │                              │
│  └──────────────┘                              │
└─────────────────────────────────────────────────┘
```

## Основные возможности

✅ **Фиатные счета** в USD, EUR, RUB, GBP  
✅ **Криптокошельки** BTC, ETH, USDT, BNB, SOL  
✅ **Переводы** между счетами  
✅ **Обмен** криптовалюты ⇄ фиат  
✅ **Аналитика** транзакций в реальном времени  
✅ **Уведомления** о всех операциях  
✅ **Мониторинг** через Prometheus и Grafana  
✅ **Трейсинг** запросов через Zipkin  

## Технологии

- 🐹 **Go 1.23** - язык программирования
- 🚀 **Fiber** - веб-фреймворк
- 🔨 **Squirrel** - SQL query builder
- ⚡ **Uber Zap** - логирование
- 🗄️ **PostgreSQL** - база данных
- 🐰 **RabbitMQ** - очередь сообщений
- 🔌 **gRPC** - межсервисное взаимодействие
- 📊 **Prometheus + Grafana** - мониторинг
- 🔍 **Zipkin** - трейсинг

