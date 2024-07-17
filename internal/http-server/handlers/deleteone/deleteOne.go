package deleteone

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func DeleteTodo(c echo.Context, store *postgresql.Storage) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID format"})
	}

	err = store.DeleteTodoItem(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Todo not found"})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting the task"})
	}

	return c.NoContent(http.StatusNoContent)
}
