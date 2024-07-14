package postgresql

import (
	"TODO_List/internal/config"
	"TODO_List/model"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Dbname,
		cfg.Database.Password,
		cfg.Database.Sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateTodoItem(todo model.Todo) (int, error) {
	const query = "INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING task_id"
	var id int
	err := s.db.QueryRow(query, todo.Title, todo.Completed).Scan(&id)
	return id, err
}

func (s *Storage) DeleteTodoItem(id int) error {
	const query = "DELETE FROM tasks WHERE task_id = $1"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *Storage) DeleteAllTodoItem() error {
	const query = "DELETE FROM tasks"
	_, err := s.db.Exec(query)
	return err
}

func (s *Storage) GetAllItems() ([]model.Todo, error) {
	rows, err := s.db.Query("SELECT task_id, title, completed FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]model.Todo, 0)

	for rows.Next() {
		var task model.Todo

		if err = rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, err
}

func (s *Storage) UpdateTodoItem(todo model.Todo, id int) error {
	const query = "UPDATE tasks SET title = $1, completed = $2 WHERE task_id = $3"
	_, err := s.db.Exec(query, todo.Title, todo.Completed, id)
	return err
}
