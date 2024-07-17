package main

import (
	"TODO_List/internal/config"
	http_server "TODO_List/internal/http-server"
	storage "TODO_List/internal/storage/postgresql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("./config/local.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	store, err := storage.New(cfg)
	if err != nil {
		log.Fatal("failed to database connect", err)
	}

	server := http_server.NewServer(store)

	e := echo.New()

	//TO DO list
	e.GET("/todos", server.GetTodos)

	e.POST("/todos", server.CreateTodo)

	e.PUT("/todos/:id", server.UpdateTodo)

	e.DELETE("/todos/:id", server.DeleteTodo)

	e.DELETE("/todos", server.DeleteAllTodos)

	//categories
	e.GET("/categories", server.GetAllCategories)

	e.GET("/categories/:id", server.GetCategoryById)

	e.POST("/categories", server.CreateCategory)

	e.PUT("/categories/:id", server.UpdateCategory)

	e.DELETE("/categories/:id", server.DeleteCategory)

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
