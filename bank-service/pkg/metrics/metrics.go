package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP metrics
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Transaction metrics
	TransactionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transactions_total",
			Help: "Total number of transactions",
		},
		[]string{"type", "status"},
	)

	TransactionAmount = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_amount",
			Help:    "Transaction amounts",
			Buckets: []float64{10, 50, 100, 500, 1000, 5000, 10000, 50000, 100000},
		},
		[]string{"currency"},
	)

	// Exchange metrics
	ExchangesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "exchanges_total",
			Help: "Total number of exchanges",
		},
		[]string{"type", "status"},
	)

	// Account metrics
	AccountsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "accounts_total",
			Help: "Total number of accounts",
		},
		[]string{"currency"},
	)

	// Wallet metrics
	WalletsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wallets_total",
			Help: "Total number of crypto wallets",
		},
		[]string{"crypto_type"},
	)
)

// InitMetrics initializes Prometheus metrics
func InitMetrics() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(TransactionsTotal)
	prometheus.MustRegister(TransactionAmount)
	prometheus.MustRegister(ExchangesTotal)
	prometheus.MustRegister(AccountsTotal)
	prometheus.MustRegister(WalletsTotal)

	// Initialize metrics with zero values to make them visible
	TransactionsTotal.WithLabelValues("transfer", "success").Add(0)
	TransactionsTotal.WithLabelValues("deposit", "success").Add(0)
	TransactionsTotal.WithLabelValues("withdraw", "success").Add(0)
	ExchangesTotal.WithLabelValues("crypto_to_fiat", "success").Add(0)
	ExchangesTotal.WithLabelValues("fiat_to_crypto", "success").Add(0)
	HttpRequestsTotal.WithLabelValues("GET", "/metrics", "200").Add(0)
}

// MetricsHandler returns Fiber handler for Prometheus metrics
func MetricsHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}
