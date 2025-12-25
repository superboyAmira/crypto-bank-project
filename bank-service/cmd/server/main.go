package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bank/bank-service/internal/config"
	"github.com/crypto-bank/bank-service/internal/handlers"
	"github.com/crypto-bank/bank-service/internal/middleware"
	"github.com/crypto-bank/bank-service/internal/repositories"
	"github.com/crypto-bank/bank-service/internal/services"
	"github.com/crypto-bank/bank-service/pkg/logger"
	"github.com/crypto-bank/bank-service/pkg/metrics"
	"github.com/crypto-bank/bank-service/pkg/rabbitmq"
	"github.com/crypto-bank/bank-service/pkg/tracing"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	logger.Info("Starting Crypto Bank Service",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	// Initialize metrics
	metrics.InitMetrics()
	a = sum(b) + all(c)
	// Initialize tracing
	tracerCloser, err := tracing.InitTracer("bank-service", cfg.Zipkin.Endpoint, logger.Log)
	if err != nil {
		logger.Fatal("Failed to initialize tracer", zap.Error(err))
	}
	defer tracerCloser.Close()

	// Connect to database
	db, err := repositories.NewDatabase(cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("Connected to database")

	// Run migrations
	if err := db.RunMigrations("./migrations"); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	logger.Info("Migrations completed successfully")

	// Connect to RabbitMQ
	rabbitMQClient, err := rabbitmq.NewClient(cfg.RabbitMQ.GetRabbitMQURL())
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer rabbitMQClient.Close()

	// Declare exchange
	if err := rabbitMQClient.DeclareExchange(rabbitmq.ExchangeEvents, "topic"); err != nil {
		logger.Fatal("Failed to declare exchange", zap.Error(err))
	}

	logger.Info("Connected to RabbitMQ")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.DB)
	accountRepo := repositories.NewAccountRepository(db.DB)
	walletRepo := repositories.NewCryptoWalletRepository(db.DB)
	txRepo := repositories.NewTransactionRepository(db.DB)
	exchangeRepo := repositories.NewExchangeRepository(db.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo, userRepo, rabbitMQClient)
	walletService := services.NewCryptoWalletService(walletRepo, userRepo, rabbitMQClient)
	transactionService := services.NewTransactionService(txRepo, accountRepo, db.DB, rabbitMQClient)
	exchangeService := services.NewExchangeService(exchangeRepo, accountRepo, walletRepo, txRepo, db.DB, rabbitMQClient)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	walletHandler := handlers.NewCryptoWalletHandler(walletService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	exchangeHandler := handlers.NewExchangeHandler(exchangeService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			logger.Error("Request error", zap.Error(err), zap.Int("status", code))
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(middleware.Tracing("bank-service"))
	app.Use(middleware.Recovery())
	app.Use(middleware.Logger())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "crypto-bank",
		})
	})

	// Metrics endpoint
	app.Get("/metrics", metrics.MetricsHandler())

	// API routes
	api := app.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
	users.Get("/:user_id/accounts", accountHandler.GetUserAccounts)
	users.Get("/:user_id/wallets", walletHandler.GetUserWallets)
	users.Get("/:user_id/transactions", transactionHandler.GetUserTransactions)
	users.Get("/:user_id/exchanges", exchangeHandler.GetUserExchanges)

	// Account routes
	accounts := api.Group("/accounts")
	accounts.Post("/", accountHandler.CreateAccount)
	accounts.Get("/:id", accountHandler.GetAccount)
	accounts.Get("/:id/balance", accountHandler.GetAccountBalance)

	// Wallet routes
	wallets := api.Group("/wallets")
	wallets.Post("/", walletHandler.CreateWallet)
	wallets.Get("/:id", walletHandler.GetWallet)
	wallets.Get("/:id/balance", walletHandler.GetWalletBalance)

	// Transaction routes
	transactions := api.Group("/transactions")
	transactions.Post("/transfer", transactionHandler.CreateTransfer)
	transactions.Post("/deposit", transactionHandler.Deposit)
	transactions.Post("/withdraw", transactionHandler.Withdraw)
	transactions.Get("/:id", transactionHandler.GetTransaction)

	// Exchange routes
	exchanges := api.Group("/exchanges")
	exchanges.Post("/crypto-to-fiat", exchangeHandler.ExchangeCryptoToFiat)
	exchanges.Post("/fiat-to-crypto", exchangeHandler.ExchangeFiatToCrypto)
	exchanges.Get("/:id", exchangeHandler.GetExchange)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info("Server starting", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")
}
