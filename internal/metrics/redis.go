package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RedisCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_hits_total",
			Help: "Total number of Redis cache hits",
		},
	)

	RedisCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_misses_total",
			Help: "Total number of Redis cache misses",
		},
	)

	RedisOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_operation_duration_seconds",
			Help:    "Duration of Redis operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(
		RedisCacheHits,
		RedisCacheMisses,
		RedisOperationDuration,
	)
}

func ObserveRedisOperation(op string, start time.Time) {
	duration := time.Since(start).Seconds()
	RedisOperationDuration.WithLabelValues(op).Observe(duration)
}