package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shutsuensha/go-tasks/internal/db"
)

type mockStore struct{}

func (m *mockStore) ListTasks(ctx context.Context) ([]db.Task, error) {
	return []db.Task{
		{
			ID:    1,
			Title: "test task",
		},
	}, nil
}

func (m *mockStore) GetTask(ctx context.Context, id int32) (db.Task, error) {
	return db.Task{
		ID:    id,
		Title: "test task",
	}, nil
}

func (m *mockStore) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	return db.Task{
		ID:    1,
		Title: arg.Title,
	}, nil
}

func TestListTasks(t *testing.T) {

	store := &mockStore{}

	handler := NewTaskHandler(store)

	r := chi.NewRouter()
	handler.Register(r)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "test task") {
		t.Fatalf("unexpected response body: %s", rr.Body.String())
	}
}