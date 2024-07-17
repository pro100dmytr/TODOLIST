package http_server

import (
	"TODO_List/internal/model"
	"TODO_List/internal/storage/postgresql"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	store *postgresql.Storage
}

func NewServer(store *postgresql.Storage) *Server {
	return &Server{store}
}
func (s *Server) GetTodos(c echo.Context) error {
	tasks, err := s.store.GetAllItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when receiving tasks"})
	}
	return c.JSON(http.StatusOK, tasks)
}

func (s *Server) CreateTodo(c echo.Context) error {
	var todo model.Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	id, err := s.store.CreateTodoItem(todo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error adding a task to the database"})
	}

	todo.ID = id
	return c.JSON(http.StatusCreated, todo)
}

func (s *Server) UpdateTodo(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	var todo model.Todo
	if err = c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	err = s.store.UpdateTodoItem(todo, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error updating a task"})
	}

	return c.JSON(http.StatusOK, todo)
}

func (s *Server) DeleteTodo(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID format"})
	}

	err = s.store.DeleteTodoItem(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Todo not found"})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting the task"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) DeleteAllTodos(c echo.Context) error {
	err := s.store.DeleteAllTodoItem()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting all tasks"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetCategoryById(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	category, err := s.store.GetCategoryByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, category)
}

func (s *Server) GetAllCategories(c echo.Context) error {
	category, err := s.store.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when receiving categories"})
	}
	return c.JSON(http.StatusOK, category)
}

func (s *Server) CreateCategory(c echo.Context) error {
	var category model.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	category.ID = 0
	id, err := s.store.CreateCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
	}
	category.ID = id
	return c.JSON(http.StatusCreated, category)
}

func (s *Server) UpdateCategory(c echo.Context) error {
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
	err = s.store.UpdateCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, category)
}

func (s *Server) DeleteCategory(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	err = s.store.DeleteCategory(id)
	if err != nil {
		if err.Error() == "category not found" {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting category"})
	}

	return c.NoContent(http.StatusNoContent)
}
