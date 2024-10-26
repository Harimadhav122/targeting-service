// metrics/metrics.go
package metrics

import (
	"github.com/go-kit/kit/metrics/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestCount   prometheus.Counter
	HttpRequestLatency prometheus.Histogram
)

func init() {
	// Initialize Prometheus metrics
	HttpRequestCount = *prometheus.NewCounterFrom(prom.CounterOpts{
		Namespace: "delivery_service",
		Subsystem: "campaigns",
		Name:      "http_request_count_total",
		Help:      "Total number of http requests",
	}, []string{"method", "code"})

	HttpRequestLatency = *prometheus.NewHistogramFrom(prom.HistogramOpts{
		Namespace: "delivery_service",
		Subsystem: "campaigns",
		Name:      "http_request_latency_seconds",
		Help:      "Request latency in seconds.",
		Buckets:   prom.DefBuckets,
	}, []string{})
}
