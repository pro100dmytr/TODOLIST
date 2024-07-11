package update

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func UpdateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Невірний ID"})
	}

	todo := new(model.Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Помилка при обробці запиту"})
	}

	query := "UPDATE tasks SET title = $1, completed = $2 WHERE task_id = $3"
	result, err := postgresql.Db.Exec(query, todo.Title, todo.Completed, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при оновленні завдання"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при отриманні кількості змінених рядків"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, todo)
}
