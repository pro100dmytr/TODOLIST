package deleteOne

import (
	"TODO_List/internal/storage/postgresql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func DeleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Невірний ID"})
	}

	query := "DELETE FROM tasks WHERE task_id = $1"
	result, err := postgresql.Db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при видаленні завдання"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при отриманні кількості змінених рядків"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}
