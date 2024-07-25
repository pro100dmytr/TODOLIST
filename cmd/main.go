package main

import (
	"TODO_List/internal/config"
	http_server "TODO_List/internal/http-server"
	"TODO_List/internal/middleware"
	storage "TODO_List/internal/storage/postgresql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("./config/prod.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	store, err := storage.New(cfg)
	if err != nil {
		log.Fatal("failed to database connect", err)
	}

	server := http_server.NewServer(store)

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

	connServer := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      e,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = connServer.ListenAndServe(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
