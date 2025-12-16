package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bank/analytics-service/internal/config"
	"github.com/crypto-bank/analytics-service/internal/service"
	"github.com/crypto-bank/analytics-service/pkg/logger"
	"github.com/crypto-bank/analytics-service/pkg/metrics"
	"github.com/crypto-bank/analytics-service/pkg/tracing"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.InitLogger(cfg.Server.Environment); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Analytics Service",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	// Initialize tracing
	tracerCloser, err := tracing.InitTracer("analytics-service", cfg.Zipkin.Endpoint, logger.Log)
	if err != nil {
		logger.Fatal("Failed to initialize tracer", zap.Error(err))
	}
	defer tracerCloser.Close()

	// Initialize analytics service
	analyticsService := service.NewAnalyticsService(logger.Log)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQ.GetRabbitMQURL())
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("Failed to open channel", zap.Error(err))
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		"bank.events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		logger.Fatal("Failed to declare exchange", zap.Error(err))
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		"analytics.queue", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		logger.Fatal("Failed to declare queue", zap.Error(err))
	}

	// Bind queue to exchange with routing keys
	routingKeys := []string{
		"transaction.#",
		"exchange.#",
		"account.#",
		"wallet.#",
	}

	for _, key := range routingKeys {
		err = ch.QueueBind(
			q.Name,        // queue name
			key,           // routing key
			"bank.events", // exchange
			false,
			nil,
		)
		if err != nil {
			logger.Fatal("Failed to bind queue", zap.Error(err), zap.String("routing_key", key))
		}
		logger.Info("Queue bound", zap.String("routing_key", key))
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Fatal("Failed to register consumer", zap.Error(err))
	}

	// Process messages in goroutine
	go func() {
		for msg := range msgs {
			logger.Debug("Received message", zap.String("routing_key", msg.RoutingKey))

			switch msg.RoutingKey {
			case "transaction.created", "transaction.completed":
				analyticsService.ProcessTransactionEvent(msg.Body)
			case "exchange.created", "exchange.completed":
				analyticsService.ProcessExchangeEvent(msg.Body)
			case "account.created":
				analyticsService.ProcessAccountEvent(msg.Body)
			case "wallet.created":
				analyticsService.ProcessWalletEvent(msg.Body)
			default:
				logger.Warn("Unknown routing key", zap.String("routing_key", msg.RoutingKey))
			}
		}
	}()

	logger.Info("Analytics service started, waiting for messages...")

	// Create HTTP server for metrics and statistics
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "analytics-service",
		})
	})

	app.Get("/metrics", metrics.MetricsHandler())

	app.Get("/api/v1/statistics", func(c *fiber.Ctx) error {
		stats := analyticsService.GetStatistics()
		return c.JSON(fiber.Map{
			"success": true,
			"data":    stats,
		})
	})

	// Start HTTP server
	go func() {
		logger.Info("HTTP server started", zap.String("port", cfg.Server.Port))
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logger.Error("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")
	app.Shutdown()
	logger.Info("Servers stopped")
}
