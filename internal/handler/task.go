package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shutsuensha/go-tasks/internal/handler/dto"
	"github.com/shutsuensha/go-tasks/internal/service"
	"github.com/shutsuensha/go-tasks/internal/validator"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) Register(r chi.Router) {
	r.Post("/tasks", h.Create)
	r.Get("/tasks", h.List)
	r.Get("/tasks/{id}", h.Get)
}

// CreateTask godoc
// @Summary Create task
// @Description create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.CreateTaskRequest true "Task data"
// @Success 201 {object} dto.TaskResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(
		r.Context(),
		req.Title,
		req.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToTaskResponse(task)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// ListTasks godoc
// @Summary List tasks
// @Description get tasks with pagination
// @Tags tasks
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} dto.TaskResponse
// @Router /tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = o
		}
	}

	tasks, err := h.service.ListTasks(
		r.Context(),
		int32(limit),
		int32(offset),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.ToTaskResponses(tasks)

	json.NewEncoder(w).Encode(resp)
}

// GetTask godoc
// @Summary Get task
// @Description get task by id
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} dto.TaskResponse
// @Failure 404 {string} string
// @Router /tasks/{id} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := dto.ToTaskResponse(task)

	json.NewEncoder(w).Encode(resp)
}