package deleteAll

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteAllTodos(c echo.Context, store *postgresql.Storage) error {
	err := store.DeleteAllTodoItem()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting all tasks"})
	}
	return c.NoContent(http.StatusNoContent)
}
