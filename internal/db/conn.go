package db

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(dsn string) (*pgxpool.Pool, error) {
    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, err
    }

    return pgxpool.NewWithConfig(context.Background(), cfg)
}