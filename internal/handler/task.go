package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shutsuensha/go-tasks/internal/db"
)

type TaskStore interface {
	ListTasks(ctx context.Context) ([]db.Task, error)
	GetTask(ctx context.Context, id int32) (db.Task, error)
	CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
}

type TaskHandler struct {
	store TaskStore
}

func NewTaskHandler(store TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) Register(r chi.Router) {
	r.Post("/tasks", h.Create)
	r.Get("/tasks", h.List)
	r.Get("/tasks/{id}", h.Get)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.store.CreateTask(r.Context(), db.CreateTaskParams{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		return
	}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {

	tasks, err := h.store.ListTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		return
	}
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.store.GetTask(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		return
	}
}