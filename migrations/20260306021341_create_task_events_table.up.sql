CREATE TABLE task_events (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);