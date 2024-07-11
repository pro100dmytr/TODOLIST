package deleteAll

import (
	"TODO_List/internal/storage/postgresql"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteAllTodos(c echo.Context) error {
	query := "DELETE FROM tasks"
	result, err := postgresql.Db.Exec(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при видаленні всіх завдань"})
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
