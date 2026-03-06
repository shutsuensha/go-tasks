-- name: CreateTask :one
INSERT INTO tasks (title, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks ORDER BY id DESC;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET done = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;