package http_server

import (
	"TODO_List/internal/model"
	"TODO_List/internal/storage/postgresql"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request"})
	}

	user, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid token"})
	}

	userID, ok := user["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid token"})
	}

	todo.UserID = int(userID) // Установка user_id

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

func (s *Server) GetCategoryTodos(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	offset := 0
	limit := 20

	if c.QueryParam("offset") != "" {
		offset, err = strconv.Atoi(c.QueryParam("offset"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid offset parameter"})
		}
	}

	if c.QueryParam("limit") != "" {
		limit, err = strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid limit parameter"})
		}
	}

	search := c.QueryParam("search")

	todos, err := s.store.GetCategoryTodos(categoryID, offset, limit, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, todos)
}

func (s *Server) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request"})
	}

	if err := validateUserData(req.Email, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to hash password"})
	}

	user := model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	id, err := s.store.CreateUser(user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return c.JSON(http.StatusConflict, model.ErrorResponse{Error: "Email already registered"})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, model.UserCreatedResponse{
		UserID: id,
	})
}

func (s *Server) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request"})
	}

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to process login"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid credentials"})
	}

	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Could not generate token"})
	}

	return c.JSON(http.StatusOK, model.TokenCreatedResponse{
		Token: "Bearer " + token,
	})
}

func (s *Server) generateJWTToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func validateUserData(email, password string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
