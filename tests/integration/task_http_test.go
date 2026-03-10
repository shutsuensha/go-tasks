package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/shutsuensha/go-tasks/internal/db"
	"github.com/shutsuensha/go-tasks/internal/handler"
	"github.com/shutsuensha/go-tasks/internal/service"
)



func TestCreateTaskHTTP(t *testing.T) {

	ctx := context.Background()

	container, connStr, err := SetupPostgres(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	err = RunMigrations(connStr)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)
	taskService := service.NewTaskService(pool, queries, nil, nil)

	taskHandler := handler.NewTaskHandler(taskService)

	r := chi.NewRouter()
	taskHandler.Register(r)

	body := map[string]string{
		"title":       "http integration test",
		"description": "test",
	}

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}