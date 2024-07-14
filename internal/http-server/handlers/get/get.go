package get

import (
	"TODO_List/internal/http-server/handlers/create"
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"

	"github.com/labstack/echo/v4"
	"net/http"
)

func GetTodos(c echo.Context, store *postgresql.Storage) error {
	rows, err := store.GetTodoItem()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Помилка при отриманні завдань"})
	}
	defer rows.Close()

	tasks := make([]model.Todo, 0)

	for rows.Next() {
		var task model.Todo
		err := rows.Scan(&task.ID, &task.Title, &task.Completed)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Помилка при скануванні рядків"})
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Помилка при ітерації по результатам запиту"})
	}

	return c.JSON(http.StatusOK, tasks)
}
