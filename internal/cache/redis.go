package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/shutsuensha/go-tasks/internal/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}