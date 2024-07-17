package getcategory

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAllCategories(c echo.Context, store *postgresql.Storage) error {
	category, err := store.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when receiving categories"})
	}
	return c.JSON(http.StatusOK, category)
}
