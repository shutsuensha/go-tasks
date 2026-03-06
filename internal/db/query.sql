-- name: CreateTask :one
INSERT INTO tasks (title, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks ORDER BY id DESC;

-- name: ListTasksPaginated :many
SELECT *
FROM tasks
ORDER BY id DESC
LIMIT $1
OFFSET $2;

-- name: CreateTaskEvent :one
INSERT INTO task_events (task_id, event_type)
VALUES ($1, $2)
RETURNING *;