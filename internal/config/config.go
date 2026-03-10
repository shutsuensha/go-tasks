package config

import (
	"github.com/spf13/viper";
    "fmt"
)

type Config struct {
	HTTPPort string
	DBUrl    string
	RedisAddr string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig() // игнорируем ошибку если файла нет

	viper.AutomaticEnv()

	cfg := &Config{
		HTTPPort: viper.GetString("HTTP_PORT"),
		DBUrl:    viper.GetString("DB_URL"),
		RedisAddr : viper.GetString("REDIS_ADDR"),
	}

	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8080"
	}

	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("DB_URL not set")
	}

	return cfg, nil
}