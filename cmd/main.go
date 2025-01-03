package main

import (
	"TODO_List/internal/config"
	http_server "TODO_List/internal/http-server"
	"TODO_List/internal/middleware"
	storage "TODO_List/internal/storage/postgresql"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("./config/prod.yaml")

	logger := config.SetupLogger(cfg.Env)

	if err != nil {
		config.Fatal(logger, "Failed to initialize database", err)
	}

	store, err := storage.New(cfg)
	if err != nil {
		config.Fatal(logger, "failed to database connect", err)
	}

	server := http_server.NewServer(store, logger)

	e := echo.New()

	authenticated := e.Group("/", middleware.JWTMiddleware)

	authenticated.GET("/todos", server.GetTodos)

	authenticated.POST("/todos", server.CreateTodo)

	authenticated.PUT("/todos/:id", server.UpdateTodo)

	authenticated.DELETE("/todos/:id", server.DeleteTodo)

	authenticated.DELETE("/todos", server.DeleteAllTodos)

	authenticated.GET("/categories", server.GetAllCategories)

	authenticated.GET("/categories/:id/todos", server.GetCategoryTodos)

	authenticated.POST("/categories", server.CreateCategory)

	authenticated.PUT("/categories/:id", server.UpdateCategory)

	authenticated.DELETE("/categories/:id", server.DeleteCategory)

	e.POST("/register", server.Register)

	e.POST("/login", server.Login)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	connServer := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      e,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	logger.Info("Server started", slog.String("address", cfg.HTTPServer.Address))

	go func() {
		if err = connServer.ListenAndServe(); err != nil {
			config.Fatal(logger, "Error starting server:", err)
		}
	}()

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = connServer.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.Any("err", err))

		return
	}

	if err = store.Close(); err != nil {
		logger.Error("Failed to close database connection", slog.Any("err", err))
	}

	logger.Info("server stopped")
}
