package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bank/exchange-service/internal/config"
	"github.com/crypto-bank/exchange-service/internal/service"
	"github.com/crypto-bank/exchange-service/pkg/logger"
	"github.com/crypto-bank/exchange-service/pkg/metrics"
	"github.com/crypto-bank/exchange-service/pkg/tracing"
	pb "github.com/crypto-bank/exchange-service/proto"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.InitLogger(cfg.Server.Environment); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Exchange Service",
		zap.String("environment", cfg.Server.Environment),
		zap.String("grpc_port", cfg.GRPC.Port),
		zap.String("http_port", cfg.Server.Port),
	)

	// Initialize tracing
	tracerCloser, err := tracing.InitTracer("exchange-service", cfg.Zipkin.Endpoint, logger.Log)
	if err != nil {
		logger.Fatal("Failed to initialize tracer", zap.Error(err))
	}
	defer tracerCloser.Close()

	// Create gRPC server with OpenTelemetry interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	exchangeService := service.NewExchangeServer(logger.Log)
	pb.RegisterExchangeServiceServer(grpcServer, exchangeService)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
		if err != nil {
			logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
		}

		logger.Info("gRPC server started", zap.String("port", cfg.GRPC.Port))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Create HTTP server for metrics and health check
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "exchange-service",
		})
	})

	app.Get("/metrics", func(c *fiber.Ctx) error {
		return c.Send([]byte("Use http://localhost:" + cfg.Server.Port + "/metrics for metrics\n"))
	})

	// HTTP server for Prometheus metrics
	go func() {
		logger.Info("Starting Prometheus metrics server", zap.String("port", "8086"))
		if err := metrics.StartMetricsServer(":8086"); err != nil {
			logger.Error("Failed to start metrics server", zap.Error(err))
		}
	}()

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
	grpcServer.GracefulStop()
	app.Shutdown()
	logger.Info("Servers stopped")
}
