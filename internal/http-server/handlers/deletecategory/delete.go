package deletecategory

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func DeleteCategory(c echo.Context, store *postgresql.Storage) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	err = store.DeleteCategory(id)
	if err != nil {
		if err.Error() == "category not found" {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting category"})
	}

	return c.NoContent(http.StatusNoContent)
}
