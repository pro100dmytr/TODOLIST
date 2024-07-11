package create

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateTodo(c echo.Context) error {
	todo := new(model.Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Помилка при обробці запиту"})
	}

	query := "INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING task_id"
	var id int
	err := postgresql.Db.QueryRow(query, todo.Title, todo.Completed).Scan(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при додаванні завдання в базу даних"})
	}

	todo.ID = id
	return c.JSON(http.StatusCreated, todo)
}
