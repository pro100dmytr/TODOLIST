package update

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func UpdateTodo(c echo.Context, store *postgresql.Storage) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	var todo model.Todo
	if err = c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	err = store.UpdateTodoItem(todo, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error updating a task"})
	}

	return c.JSON(http.StatusOK, todo)
}
