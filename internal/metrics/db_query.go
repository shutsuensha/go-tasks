package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var DBQueryDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Duration of database queries",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"query"},
)

func init() {
	prometheus.MustRegister(DBQueryDuration)
}

func ObserveQuery(query string, start time.Time) {
	duration := time.Since(start).Seconds()
	DBQueryDuration.WithLabelValues(query).Observe(duration)
}