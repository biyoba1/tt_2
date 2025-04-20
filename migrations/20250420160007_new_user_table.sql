-- +goose Up
CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    task_identification VARCHAR(50) NOT NULL,
    task VARCHAR(50) NOT NULL,
    result TEXT NOT NULL,
    status VARCHAR(15) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS task;
