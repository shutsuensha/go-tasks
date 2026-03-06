package metrics

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartDBCollector(pool *pgxpool.Pool) {

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			CollectDBStats(pool)
		}
	}()
}