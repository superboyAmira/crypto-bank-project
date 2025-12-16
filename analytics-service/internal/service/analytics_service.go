package service

import (
	"encoding/json"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	// Analytics metrics
	TransactionsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_transactions_processed_total",
			Help: "Total number of transactions processed by analytics",
		},
		[]string{"type", "status"},
	)

	ExchangesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_exchanges_processed_total",
			Help: "Total number of exchanges processed by analytics",
		},
		[]string{"type", "status"},
	)

	AccountsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_accounts_created_total",
			Help: "Total number of accounts created",
		},
	)

	WalletsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_wallets_created_total",
			Help: "Total number of wallets created",
		},
	)

	TransactionVolume = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "analytics_transaction_volume",
			Help: "Total transaction volume",
		},
		[]string{"currency"},
	)
)

func init() {
	prometheus.MustRegister(TransactionsProcessed)
	prometheus.MustRegister(ExchangesProcessed)
	prometheus.MustRegister(AccountsCreated)
	prometheus.MustRegister(WalletsCreated)
	prometheus.MustRegister(TransactionVolume)
	
	// Initialize metrics with zero values to make them visible
	TransactionsProcessed.WithLabelValues("TRANSFER", "completed").Add(0)
	TransactionsProcessed.WithLabelValues("DEPOSIT", "completed").Add(0)
	TransactionsProcessed.WithLabelValues("WITHDRAW", "completed").Add(0)
	AccountsCreated.Add(0)
	WalletsCreated.Add(0)
}

type AnalyticsService struct {
	logger *zap.Logger
	stats  *Statistics
	mu     sync.RWMutex
}

type Statistics struct {
	TotalTransactions   int64              `json:"total_transactions"`
	TotalExchanges      int64              `json:"total_exchanges"`
	TotalAccounts       int64              `json:"total_accounts"`
	TotalWallets        int64              `json:"total_wallets"`
	TransactionsByType  map[string]int64   `json:"transactions_by_type"`
	ExchangesByType     map[string]int64   `json:"exchanges_by_type"`
	VolumesByCurrency   map[string]float64 `json:"volumes_by_currency"`
}

func NewAnalyticsService(logger *zap.Logger) *AnalyticsService {
	return &AnalyticsService{
		logger: logger,
		stats: &Statistics{
			TransactionsByType: make(map[string]int64),
			ExchangesByType:    make(map[string]int64),
			VolumesByCurrency:  make(map[string]float64),
		},
	}
}

// ProcessTransactionEvent processes transaction events
func (s *AnalyticsService) ProcessTransactionEvent(body []byte) error {
	var event struct {
		TransactionID string  `json:"transaction_id"`
		UserID        string  `json:"user_id"`
		Type          string  `json:"type"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		Status        string  `json:"status"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal transaction event", zap.Error(err))
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.TotalTransactions++
	s.stats.TransactionsByType[event.Type]++
	s.stats.VolumesByCurrency[event.Currency] += event.Amount

	// Update metrics
	TransactionsProcessed.WithLabelValues(event.Type, event.Status).Inc()
	TransactionVolume.WithLabelValues(event.Currency).Set(s.stats.VolumesByCurrency[event.Currency])

	s.logger.Info("Transaction event processed",
		zap.String("transaction_id", event.TransactionID),
		zap.String("type", event.Type),
		zap.Float64("amount", event.Amount),
		zap.String("currency", event.Currency),
	)

	return nil
}

// ProcessExchangeEvent processes exchange events
func (s *AnalyticsService) ProcessExchangeEvent(body []byte) error {
	var event struct {
		ExchangeID   string  `json:"exchange_id"`
		UserID       string  `json:"user_id"`
		FromCurrency string  `json:"from_currency"`
		ToCurrency   string  `json:"to_currency"`
		FromAmount   float64 `json:"from_amount"`
		ToAmount     float64 `json:"to_amount"`
		Status       string  `json:"status"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal exchange event", zap.Error(err))
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.TotalExchanges++
	exchangeType := event.FromCurrency + "_to_" + event.ToCurrency
	s.stats.ExchangesByType[exchangeType]++

	// Update metrics
	ExchangesProcessed.WithLabelValues(exchangeType, event.Status).Inc()

	s.logger.Info("Exchange event processed",
		zap.String("exchange_id", event.ExchangeID),
		zap.String("from", event.FromCurrency),
		zap.String("to", event.ToCurrency),
		zap.Float64("amount", event.FromAmount),
	)

	return nil
}

// ProcessAccountEvent processes account creation events
func (s *AnalyticsService) ProcessAccountEvent(body []byte) error {
	var event struct {
		AccountID string `json:"account_id"`
		UserID    string `json:"user_id"`
		Currency  string `json:"currency"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal account event", zap.Error(err))
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.TotalAccounts++
	AccountsCreated.Inc()

	s.logger.Info("Account event processed",
		zap.String("account_id", event.AccountID),
		zap.String("currency", event.Currency),
	)

	return nil
}

// ProcessWalletEvent processes wallet creation events
func (s *AnalyticsService) ProcessWalletEvent(body []byte) error {
	var event struct {
		WalletID   string `json:"wallet_id"`
		UserID     string `json:"user_id"`
		CryptoType string `json:"crypto_type"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal wallet event", zap.Error(err))
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.TotalWallets++
	WalletsCreated.Inc()

	s.logger.Info("Wallet event processed",
		zap.String("wallet_id", event.WalletID),
		zap.String("crypto_type", event.CryptoType),
	)

	return nil
}

// GetStatistics returns current statistics
func (s *AnalyticsService) GetStatistics() *Statistics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy of statistics
	statsCopy := &Statistics{
		TotalTransactions:  s.stats.TotalTransactions,
		TotalExchanges:     s.stats.TotalExchanges,
		TotalAccounts:      s.stats.TotalAccounts,
		TotalWallets:       s.stats.TotalWallets,
		TransactionsByType: make(map[string]int64),
		ExchangesByType:    make(map[string]int64),
		VolumesByCurrency:  make(map[string]float64),
	}

	for k, v := range s.stats.TransactionsByType {
		statsCopy.TransactionsByType[k] = v
	}
	for k, v := range s.stats.ExchangesByType {
		statsCopy.ExchangesByType[k] = v
	}
	for k, v := range s.stats.VolumesByCurrency {
		statsCopy.VolumesByCurrency[k] = v
	}

	return statsCopy
}

