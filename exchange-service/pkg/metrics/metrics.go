package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// gRPC metrics
	GrpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	// Exchange metrics
	ExchangesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "exchanges_total",
			Help: "Total number of currency exchanges",
		},
		[]string{"from_currency", "to_currency", "status"},
	)
)

func init() {
	prometheus.MustRegister(GrpcRequestsTotal)
	prometheus.MustRegister(ExchangesTotal)

	// Initialize metrics with zero values to make them visible
	GrpcRequestsTotal.WithLabelValues("GetExchangeRate", "success").Add(0)
	GrpcRequestsTotal.WithLabelValues("GetAllRates", "success").Add(0)
	GrpcRequestsTotal.WithLabelValues("UpdateRate", "success").Add(0)
	ExchangesTotal.WithLabelValues("BTC", "USD", "success").Add(0)
	ExchangesTotal.WithLabelValues("ETH", "USD", "success").Add(0)
}

// Handler returns the Prometheus HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}

// StartMetricsServer starts the HTTP server for metrics
func StartMetricsServer(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}
