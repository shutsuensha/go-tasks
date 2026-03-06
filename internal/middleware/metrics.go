package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	httpInflightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_inflight_requests",
			Help: "Current number of inflight requests",
		},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpInflightRequests,
	)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		httpInflightRequests.Inc()
		defer httpInflightRequests.Dec()

		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		status := strconv.Itoa(rw.status)

		httpRequestsTotal.WithLabelValues(path, method, status).Inc()

		httpRequestDuration.WithLabelValues(path, method).Observe(duration)
	})
}