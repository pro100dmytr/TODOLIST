package createcategory

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateCategory(c echo.Context, store *postgresql.Storage) error {
	var category model.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	category.ID = 0
	id, err := store.CreateCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
	}
	category.ID = id
	return c.JSON(http.StatusCreated, category)
}
