-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    task_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL,
    category_id INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd