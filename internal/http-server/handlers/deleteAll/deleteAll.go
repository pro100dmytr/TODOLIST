package deleteAll

import (
	"TODO_List/internal/http-server/handlers/create"
	"TODO_List/internal/storage/postgresql"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteAllTodos(c echo.Context, store *postgresql.Storage) error {
	result, err := store.DeleteAllTodoItem()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Error when deleting all tasks"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Error when getting the number of changed rows"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}
