.PHONY: help build up down logs clean proto test

build:
	docker-compose build

up:
	docker-compose up -d

jenkins-up:
	docker-compose -f docker-compose.jenkins.yml up -d
	docker exec crypto_bank_jenkins cat /var/jenkins_home/secrets/initialAdminPassword
# plugins:
# docker pipline
# pipeline
# Git server
# Go

down:
	docker-compose down

logs:
	docker-compose logs -f

logs-bank:
	docker-compose logs -f bank-service

logs-exchange:
	docker-compose logs -f exchange-service

logs-analytics:
	docker-compose logs -f analytics-service

logs-notification:
	docker-compose logs -f notification-service

clean:
	docker-compose down -v --rmi all

restart:
	docker-compose restart

ps:
	docker-compose ps

proto:
	cd exchange-service && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/exchange.proto


migrate-up:
	cd bank-service && goose -dir migrations postgres "host=localhost port=5432 user=postgres password=1234 dbname=crypto_bank sslmode=disable" up

migrate-down:
	cd bank-service && goose -dir migrations postgres "host=localhost port=5432 user=postgres password=1234 dbname=crypto_bank sslmode=disable" down

test:
	cd bank-service && go test ./... -v
	cd exchange-service && go test ./... -v
	cd analytics-service && go test ./... -v
	cd notification-service && go test ./... -v

open-grafana:
	@echo "Opening Grafana at http://localhost:3000"
	@echo "Username: admin, Password: admin"
	@open http://localhost:3000 || xdg-open http://localhost:3000

open-prometheus:
	@echo "Opening Prometheus at http://localhost:9091"
	@open http://localhost:9091 || xdg-open http://localhost:9091

open-zipkin:
	@echo "Opening Zipkin at http://localhost:9411"
	@open http://localhost:9411 || xdg-open http://localhost:9411

open-rabbitmq:
	@echo "Opening RabbitMQ Management at http://localhost:15672"
	@echo "Username: guest, Password: guest"
	@open http://localhost:15672 || xdg-open http://localhost:15672

api-health:
	@echo "Bank Service:"
	@curl -s http://localhost:8080/health | jq
	@echo "\nExchange Service:"
	@curl -s http://localhost:8085/health | jq
	@echo "\nAnalytics Service:"
	@curl -s http://localhost:8082/health | jq
	@echo "\nNotification Service:"
	@curl -s http://localhost:8083/health | jq

api-stats:
	@curl -s http://localhost:8082/api/v1/statistics | jq

api-notifications:
	@curl -s http://localhost:8083/api/v1/notifications | jq

# Windows build commands (for local development without Docker)
build-windows:
	@echo "Building bank-service for Windows..."
	cd bank-service && set CGO_ENABLED=0 && set GOOS=windows && go build -ldflags="-w -s" -o main.exe ./cmd/server
	@echo "Building exchange-service for Windows..."
	cd exchange-service && set CGO_ENABLED=0 && set GOOS=windows && go build -ldflags="-w -s" -o main.exe ./cmd/server
	@echo "Building analytics-service for Windows..."
	cd analytics-service && set CGO_ENABLED=0 && set GOOS=windows && go build -ldflags="-w -s" -o main.exe ./cmd/server
	@echo "Building notification-service for Windows..."
	cd notification-service && set CGO_ENABLED=0 && set GOOS=windows && go build -ldflags="-w -s" -o main.exe ./cmd/server

build-linux:
	@echo "Building bank-service for Linux..."
	cd bank-service && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server
	@echo "Building exchange-service for Linux..."
	cd exchange-service && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server
	@echo "Building analytics-service for Linux..."
	cd analytics-service && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server
	@echo "Building notification-service for Linux..."
	cd notification-service && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server
