package update

import (
	"TODO_List/internal/http-server/handlers/create"
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func UpdateTodo(c echo.Context, store *postgresql.Storage) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, create.ErrorResponse{Error: "Invalid ID"})
	}

	var todo model.Todo
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, create.ErrorResponse{Error: "Error while processing a request"})
	}

	result, err := store.UpdateTodoItem(todo, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Error updating a task"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, create.ErrorResponse{Error: "Error getting the number of changed rows"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, todo)
}
