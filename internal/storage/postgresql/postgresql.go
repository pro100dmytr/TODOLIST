package postgresql

import (
	"TODO_List/internal/config"
	"TODO_List/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pressly/goose/v3"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Close() error {
	return s.db.Close()
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
		slog.Error("Failed to open database connection", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		slog.Error("Failed to ping database", err)
		return nil, err
	}

	err = goose.Up(db, "db\\migrations")
	if err != nil {
		slog.Error("Failed to run migrations", err)
		return nil, err
	}

	slog.Info("Connected to database and initialized is successfully")
	return &Storage{db: db}, nil
}

func (s *Storage) CreateTodoItem(ctx context.Context, todo model.Todo) (int, error) {
	const query = "INSERT INTO tasks (title, completed, category_id, user_id) VALUES ($1, $2, $3, $4) RETURNING task_id"
	var id int
	err := s.db.QueryRowContext(ctx, query, todo.Title, todo.Completed, todo.CategoryID, todo.UserID).Scan(&id)
	return id, err
}

func (s *Storage) DeleteTodoItem(ctx context.Context, id int) error {
	const query = "DELETE FROM tasks WHERE task_id = $1"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *Storage) DeleteAllTodoItem(ctx context.Context) error {
	const query = "DELETE FROM tasks"
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *Storage) GetAllItems(ctx context.Context) ([]model.Todo, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT task_id, title, completed, category_id FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]model.Todo, 0)

	for rows.Next() {
		var task model.Todo

		if err = rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CategoryID); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) UpdateTodoItem(ctx context.Context, todo model.Todo, id int) error {
	const query = "UPDATE tasks SET title = $1, completed = $2, category_id = $3 WHERE task_id = $4"
	_, err := s.db.ExecContext(ctx, query, todo.Title, todo.Completed, todo.CategoryID, id)
	return err
}

func (s *Storage) CreateCategory(ctx context.Context, category model.Category) (int, error) {
	const query = "INSERT INTO categories (category) VALUES ($1) RETURNING id"
	var id int
	err := s.db.QueryRowContext(ctx, query, category.Category).Scan(&id)
	return id, err
}

func (s *Storage) UpdateCategory(ctx context.Context, category model.Category) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		"UPDATE categories SET category = $1 WHERE id = $2",
		category.Category,
		category.ID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE tasks SET category_id = $1 WHERE task_id = $2",
		category.ID,
		category.ID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) DeleteCategory(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM tasks WHERE category_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	const query = "SELECT id, category FROM categories"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err = rows.Scan(&c.ID, &c.Category); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *Storage) GetCategoryTodos(ctx context.Context, categoryID, offset, limit int, search string) ([]model.Todo, error) {
	query := `
    SELECT task_id, title, completed, category_id
    FROM tasks
    WHERE category_id = $1
`
	args := []interface{}{categoryID}

	if search != "" {
		query += " AND title LIKE $2"
		args = append(args, "%"+search+"%")
	}

	query += `
    ORDER BY task_id
    LIMIT $%d
    OFFSET $%d 
`

	args = append(args, limit, offset)

	query = fmt.Sprintf(query, len(args)-1, len(args))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Todo

	for rows.Next() {
		var task model.Todo
		err = rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CategoryID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) CreateUser(ctx context.Context, user model.User) (int, error) {
	const query = "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING user_id"
	var id int
	err := s.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&id)
	return id, err
}

var ErrUserNotFound = errors.New("user not found")

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	query := `
		SELECT user_id, email, password 
		FROM users
		WHERE email = $1
	`

	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrUserNotFound
	}

	return user, err
}
