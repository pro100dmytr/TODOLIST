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

	r := e.Group("/", middleware.JWTMiddleware)
	{
		r.GET("/todos", server.GetTodos)

		r.POST("/todoss", server.CreateTodo)

		r.PUT("/todos/:id", server.UpdateTodo)

		r.DELETE("/todos/:id", server.DeleteTodo)

		r.DELETE("/todos", server.DeleteAllTodos)

		r.GET("/categories", server.GetAllCategories)

		r.GET("/categories/:id/todos", server.GetCategoryTodos)

		r.POST("/categories", server.CreateCategory)

		r.PUT("/categories/:id", server.UpdateCategory)

		r.DELETE("/categories/:id", server.DeleteCategory)

	}

	//authentication and authorization
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
