package updatecategory

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func UpdateCategory(c echo.Context, store *postgresql.Storage) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	var category model.Category
	if err = c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	category.ID = id
	err = store.UpdateCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, category)
}
