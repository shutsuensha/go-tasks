package metrics

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dbTotalConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_total_connections",
			Help: "Total number of connections in the pool",
		},
	)

	dbIdleConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_idle_connections",
			Help: "Idle connections in pool",
		},
	)

	dbAcquiredConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_acquired_connections",
			Help: "Connections currently in use",
		},
	)

	dbConstructingConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_constructing_connections",
			Help: "Connections being created",
		},
	)

	dbMaxConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_max_connections",
			Help: "Maximum pool size",
		},
	)
)

func init() {
	prometheus.MustRegister(
		dbTotalConnections,
		dbIdleConnections,
		dbAcquiredConnections,
		dbConstructingConnections,
		dbMaxConnections,
	)
}

func CollectDBStats(pool *pgxpool.Pool) {

	stats := pool.Stat()

	dbTotalConnections.Set(float64(stats.TotalConns()))
	dbIdleConnections.Set(float64(stats.IdleConns()))
	dbAcquiredConnections.Set(float64(stats.AcquiredConns()))
	dbConstructingConnections.Set(float64(stats.ConstructingConns()))
	dbMaxConnections.Set(float64(stats.MaxConns()))
}