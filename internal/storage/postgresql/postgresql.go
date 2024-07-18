package postgresql

import (
	"TODO_List/internal/config"
	"TODO_List/internal/model"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
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

	//migration
	err = goose.Up(db, "db\\migrations")
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateTodoItem(todo model.Todo) (int, error) {
	const query = "INSERT INTO tasks (title, completed, category_id) VALUES ($1, $2, $3) RETURNING task_id"
	var id int
	err := s.db.QueryRow(query, todo.Title, todo.Completed, todo.CategoryID).Scan(&id)
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
	rows, err := s.db.Query("SELECT task_id, title, completed, category_id FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks, err := scanToDo(rows, func(rows *sql.Rows, task *model.Todo) error {
		return rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CategoryID)
	})
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, err
}

func (s *Storage) UpdateTodoItem(todo model.Todo, id int) error {
	const query = "UPDATE tasks SET title = $1, completed = $2, category_id = $3 WHERE task_id = $4"
	_, err := s.db.Exec(query, todo.Title, todo.Completed, todo.CategoryID, id)
	return err
}

func (s *Storage) CreateCategory(category model.Category) (id2 int, err error) {
	const query = "INSERT INTO categories (category) VALUES ($1) RETURNING id"
	return id2, s.db.QueryRow(query, category.Category).Scan(&id2)
}

func (s *Storage) UpdateCategory(category model.Category) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE categories SET category = $1 WHERE id = $2",
		category.Category,
		category.ID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE tasks SET category_id = $1 WHERE task_id = $2",
		category.Category,
		category.ID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) DeleteCategory(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM tasks WHERE category_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetAllCategories() ([]model.Category, error) {
	const query = "SELECT id, category FROM categories"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories, err := scanToDo(rows, func(rows *sql.Rows, category *model.Category) error {
		return rows.Scan(&category.ID, &category.Category)
	})
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *Storage) GetCategoryByID(id int) (model.Category, error) {
	var category model.Category

	err := s.db.QueryRow("SELECT id, category FROM categories WHERE id = $1", id).
		Scan(&category.ID, &category.Category)
	if err != nil {
		return category, err
	}

	rows, err := s.db.Query("SELECT task_id, title, completed, category_id FROM tasks WHERE category_id = $1", id)
	if err != nil {
		return category, err
	}
	defer rows.Close()

	tasks, err := scanToDo(rows, func(rows *sql.Rows, task *model.Todo) error {
		return rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CategoryID)
	})
	if err != nil {
		return category, err
	}

	category.Tasks = tasks
	return category, nil
}

func scanToDo[T any](rows *sql.Rows, scanFunc func(*sql.Rows, *T) error) ([]T, error) {
	items := make([]T, 0)
	for rows.Next() {
		var item T
		if err := scanFunc(rows, &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
