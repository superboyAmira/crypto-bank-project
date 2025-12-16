package service

import (
	"context"
	"sync"
	"time"

	"github.com/crypto-bank/exchange-service/pkg/metrics"
	pb "github.com/crypto-bank/exchange-service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExchangeServer struct {
	pb.UnimplementedExchangeServiceServer
	rates  map[string]float64
	mu     sync.RWMutex
	logger *zap.Logger
}

func NewExchangeServer(logger *zap.Logger) *ExchangeServer {
	server := &ExchangeServer{
		rates:  make(map[string]float64),
		logger: logger,
	}

	// Initialize with default rates
	server.initializeRates()

	return server
}

func (s *ExchangeServer) initializeRates() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Crypto to USD
	s.rates["BTC-USD"] = 43500.00
	s.rates["ETH-USD"] = 2280.50
	s.rates["USDT-USD"] = 1.00
	s.rates["BNB-USD"] = 315.75
	s.rates["SOL-USD"] = 98.30

	// USD to Crypto
	s.rates["USD-BTC"] = 0.000023
	s.rates["USD-ETH"] = 0.000438
	s.rates["USD-USDT"] = 1.00
	s.rates["USD-BNB"] = 0.003167
	s.rates["USD-SOL"] = 0.010173

	// Fiat conversions
	s.rates["USD-EUR"] = 0.92
	s.rates["USD-RUB"] = 92.50
	s.rates["USD-GBP"] = 0.79
	s.rates["EUR-USD"] = 1.09
	s.rates["RUB-USD"] = 0.0108
	s.rates["GBP-USD"] = 1.27

	// Crypto to other fiat
	s.rates["BTC-EUR"] = 40020.00
	s.rates["ETH-EUR"] = 2097.66
	s.rates["BTC-RUB"] = 4023750.00
	s.rates["ETH-RUB"] = 210941.25

	s.logger.Info("Exchange rates initialized", zap.Int("count", len(s.rates)))
}

func (s *ExchangeServer) GetExchangeRate(ctx context.Context, req *pb.ExchangeRateRequest) (*pb.ExchangeRateResponse, error) {
	s.logger.Info("GetExchangeRate called",
		zap.String("from", req.FromCurrency),
		zap.String("to", req.ToCurrency),
	)

	s.mu.RLock()
	defer s.mu.RUnlock()

	key := req.FromCurrency + "-" + req.ToCurrency
	rate, exists := s.rates[key]

	if !exists {
		s.logger.Warn("Exchange rate not found",
			zap.String("from", req.FromCurrency),
			zap.String("to", req.ToCurrency),
		)
		metrics.GrpcRequestsTotal.WithLabelValues("GetExchangeRate", "error").Inc()
		return nil, status.Errorf(codes.NotFound, "exchange rate not found for %s to %s", req.FromCurrency, req.ToCurrency)
	}

	metrics.GrpcRequestsTotal.WithLabelValues("GetExchangeRate", "success").Inc()
	metrics.ExchangesTotal.WithLabelValues(req.FromCurrency, req.ToCurrency, "success").Inc()

	return &pb.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate,
		Timestamp:    time.Now().Unix(),
	}, nil
}

func (s *ExchangeServer) GetAllRates(ctx context.Context, req *pb.Empty) (*pb.AllRatesResponse, error) {
	s.logger.Info("GetAllRates called")

	s.mu.RLock()
	defer s.mu.RUnlock()

	var rates []*pb.ExchangeRateResponse
	timestamp := time.Now().Unix()

	for key, rate := range s.rates {
		// Parse key (e.g., "BTC-USD" -> "BTC" and "USD")
		fromCurrency := key[:len(key)-4]
		toCurrency := key[len(key)-3:]

		rates = append(rates, &pb.ExchangeRateResponse{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         rate,
			Timestamp:    timestamp,
		})
	}

	metrics.GrpcRequestsTotal.WithLabelValues("GetAllRates", "success").Inc()

	return &pb.AllRatesResponse{
		Rates: rates,
	}, nil
}

func (s *ExchangeServer) UpdateRate(ctx context.Context, req *pb.UpdateRateRequest) (*pb.UpdateRateResponse, error) {
	s.logger.Info("UpdateRate called",
		zap.String("from", req.FromCurrency),
		zap.String("to", req.ToCurrency),
		zap.Float64("rate", req.Rate),
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	key := req.FromCurrency + "-" + req.ToCurrency
	s.rates[key] = req.Rate

	// Also update inverse rate
	inverseKey := req.ToCurrency + "-" + req.FromCurrency
	if req.Rate != 0 {
		s.rates[inverseKey] = 1.0 / req.Rate
	}

	s.logger.Info("Exchange rate updated",
		zap.String("key", key),
		zap.Float64("rate", req.Rate),
	)

	metrics.GrpcRequestsTotal.WithLabelValues("UpdateRate", "success").Inc()

	return &pb.UpdateRateResponse{
		Success: true,
		Message: "Rate updated successfully",
	}, nil
}
