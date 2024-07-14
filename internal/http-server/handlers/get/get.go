package get

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"

	"github.com/labstack/echo/v4"
	"net/http"
)

func GetTodos(c echo.Context, store *postgresql.Storage) error {
	tasks, err := store.GetAllItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when receiving tasks"})
	}
	return c.JSON(http.StatusOK, tasks)
}
