package deleteOne

import (
	"TODO_List/internal/http-server/handlers/create"
	"TODO_List/internal/storage/postgresql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func DeleteTodo(c echo.Context, store *postgresql.Storage) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, create.ErrorResponse{Error: "Invalid ID"})
	}

	result, err := store.DeleteTodoItem(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Error when deleting all tasks"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Error getting the number of changed rows"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}
