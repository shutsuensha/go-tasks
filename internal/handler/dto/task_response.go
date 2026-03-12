package dto

import (
	"time"

	"github.com/shutsuensha/go-tasks/internal/db"
)

type TaskResponse struct {
	ID          int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToTaskResponse(task db.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt.Time,
	}
}

func ToTaskResponses(tasks []db.Task) []TaskResponse {
	resp := make([]TaskResponse, 0, len(tasks))

	for _, t := range tasks {
		resp = append(resp, ToTaskResponse(t))
	}

	return resp
}