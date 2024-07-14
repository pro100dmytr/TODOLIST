package deleteOne

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func DeleteTodo(c echo.Context, store *postgresql.Storage) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	err = store.DeleteTodoItem(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting all tasks"})
	}

	return c.NoContent(http.StatusNoContent)
}
