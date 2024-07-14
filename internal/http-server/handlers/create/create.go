package create

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateTodo(c echo.Context, store *postgresql.Storage) error {
	var todo model.Todo
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	id, err := store.CreateTodoItem(todo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error adding a task to the database"})
	}

	todo.ID = id
	return c.JSON(http.StatusCreated, todo)
}
