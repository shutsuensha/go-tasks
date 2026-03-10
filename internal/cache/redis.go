package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/shutsuensha/go-tasks/internal/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	return rdb
}