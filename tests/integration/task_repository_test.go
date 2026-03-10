package integration

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/shutsuensha/go-tasks/internal/db"
)

func TestCreateTask(t *testing.T) {

	ctx := context.Background()

	container, connStr, err := SetupPostgres(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// запускаем migrations
	err = RunMigrations(connStr)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	task, err := queries.CreateTask(ctx, db.CreateTaskParams{
		Title:       "integration test",
		Description: "desct",
	})

	require.NoError(t, err)
	require.Equal(t, "integration test", task.Title)
}