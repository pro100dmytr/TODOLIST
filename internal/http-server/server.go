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
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	store  *postgresql.Storage
	logger *slog.Logger
}

func NewServer(store *postgresql.Storage, logger *slog.Logger) *Server {
	return &Server{store: store, logger: logger}
}

func (s *Server) GetTodos(c echo.Context) error {
	ctx := c.Request().Context()
	tasks, err := s.store.GetAllItems(ctx)
	if err != nil {
		s.logger.Error("Error when receiving tasks", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when receiving tasks"})
	}

	return c.JSON(http.StatusOK, tasks)
}

func (s *Server) CreateTodo(c echo.Context) error {
	var todo model.Todo
	if err := c.Bind(&todo); err != nil {
		s.logger.Error("Invalid request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request"})
	}

	user, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		s.logger.Error("Invalid token")
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid token"})
	}

	userID, ok := user["user_id"].(float64)
	if !ok {
		s.logger.Error("Invalid token")
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid token"})
	}

	todo.UserID = int(userID)

	ctx := c.Request().Context()
	id, err := s.store.CreateTodoItem(ctx, todo)
	if err != nil {
		s.logger.Error("Error adding a task to the database", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error adding a task to the database"})
	}

	todo.ID = id

	s.logger.Info("Task created", slog.String("id", strconv.Itoa(id)))

	return c.JSON(http.StatusCreated, todo)
}

func (s *Server) UpdateTodo(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.logger.Error("Invalid ID", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID"})
	}

	var todo model.Todo
	if err = c.Bind(&todo); err != nil {
		s.logger.Error("Error while processing a request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	ctx := c.Request().Context()

	err = s.store.UpdateTodoItem(ctx, todo, id)
	if err != nil {
		s.logger.Error("Error updating a task", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error updating a task"})
	}

	s.logger.Info("Task updated", slog.String("id", idStr))

	return c.JSON(http.StatusOK, todo)
}

func (s *Server) DeleteTodo(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.logger.Error("Invalid ID format", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID format"})
	}

	ctx := c.Request().Context()

	err = s.store.DeleteTodoItem(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("Todo task not found", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Todo task not found"})
		}

		s.logger.Error("Error when deleting the task", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting the task"})
	}

	s.logger.Info("Task deleted", slog.String("id", idStr))
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) DeleteAllTodos(c echo.Context) error {
	ctx := c.Request().Context()

	err := s.store.DeleteAllTodoItem(ctx)
	if err != nil {
		s.logger.Error("Error when deleting all tasks", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting all tasks"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetAllCategories(c echo.Context) error {
	ctx := c.Request().Context()

	category, err := s.store.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error("Error when receiving categories", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when receiving categories"})
	}

	return c.JSON(http.StatusOK, category)
}

func (s *Server) CreateCategory(c echo.Context) error {
	var category model.Category
	if err := c.Bind(&category); err != nil {
		s.logger.Error("Error while processing a request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request"})
	}

	category.ID = 0
	ctx := c.Request().Context()

	id, err := s.store.CreateCategory(ctx, category)
	if err != nil {
		s.logger.Error("Error creating a category", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error creating a category"})
	}

	category.ID = id

	s.logger.Info("Category updated", slog.String("id", strconv.Itoa(id)))

	return c.JSON(http.StatusCreated, category)
}

func (s *Server) UpdateCategory(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.logger.Error("Invalid ID category", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID category"})
	}

	var category model.Category
	if err = c.Bind(&category); err != nil {
		s.logger.Error("Error while processing a request category", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Error while processing a request category"})
	}

	category.ID = id
	ctx := c.Request().Context()

	err = s.store.UpdateCategory(ctx, category)
	if err != nil {
		s.logger.Error("Error updating a category", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error updating a category"})
	}

	s.logger.Info("Category updated", slog.String("id", idStr))
	return c.JSON(http.StatusOK, category)
}

func (s *Server) DeleteCategory(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.logger.Error("Invalid ID category", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID category"})
	}

	ctx := c.Request().Context()

	err = s.store.DeleteCategory(ctx, id)
	if err != nil {
		if err.Error() == "category not found" {
			s.logger.Error("Category not found", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
		}

		s.logger.Error("Error when deleting category", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Error when deleting category"})
	}

	s.logger.Info("Category deleted", slog.String("id", idStr))
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetCategoryTodos(c echo.Context) error {
	idStr := strings.TrimPrefix(c.Param("id"), ":")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		s.logger.Error("Invalid ID category", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid ID category"})
	}

	offset := 0
	limit := 20

	if c.QueryParam("offset") != "" {
		offset, err = strconv.Atoi(c.QueryParam("offset"))
		if err != nil {
			s.logger.Error("Invalid offset parameter", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid offset parameter"})
		}
	}

	if c.QueryParam("limit") != "" {
		limit, err = strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			s.logger.Error("Invalid limit parameter", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid limit parameter"})
		}
	}

	search := c.QueryParam("search")
	ctx := c.Request().Context()

	todos, err := s.store.GetCategoryTodos(ctx, categoryID, offset, limit, search)
	if err != nil {
		s.logger.Error("Invalid limit parameter", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid limit parameter"})
	}

	return c.JSON(http.StatusOK, todos)
}

func (s *Server) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("Invalid request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request"})
	}

	if err := validateUserData(req.Email, req.Password); err != nil {
		s.logger.Error("Invalid user validation", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid user validation"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to hash password"})
	}

	user := model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	ctx := c.Request().Context()

	id, err := s.store.CreateUser(ctx, user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			s.logger.Error("Email already registered", slog.Any("error", err))
			return c.JSON(http.StatusConflict, model.ErrorResponse{Error: "Email already registered"})
		}

		s.logger.Error("Failed to create user", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create user"})
	}

	s.logger.Info("User created", slog.String("email", user.Email))
	return c.JSON(http.StatusCreated, model.UserCreatedResponse{
		UserID: id,
	})
}

func (s *Server) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("Invalid request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request"})
	}

	ctx := c.Request().Context()

	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			s.logger.Error("Invalid credentials", slog.Any("error", err))
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid credentials"})
		}

		s.logger.Error("Failed to process login", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to process login"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Error("Invalid credentials", slog.Any("error", err))
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid credentials"})
	}

	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		s.logger.Error("Could not generate token", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Could not generate token"})
	}

	s.logger.Info("User log ginned", slog.String("email", user.Email))
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
