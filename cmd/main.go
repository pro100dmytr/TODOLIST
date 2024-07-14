package main

import (
	"TODO_List/internal/config"
	"TODO_List/internal/http-server/handlers/create"
	"TODO_List/internal/http-server/handlers/deleteAll"
	"TODO_List/internal/http-server/handlers/deleteOne"
	"TODO_List/internal/http-server/handlers/get"
	"TODO_List/internal/http-server/handlers/update"
	storage "TODO_List/internal/storage/postgresql"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
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

	e := echo.New()

	e.GET("/todos", func(c echo.Context) error {
		return get.GetTodos(c, store)
	})
	e.POST("/todos", func(c echo.Context) error {
		return create.CreateTodo(c, store)
	})
	e.PUT("/todos/:id", func(c echo.Context) error {
		return update.UpdateTodo(c, store)
	})
	e.DELETE("/todos/:id", func(c echo.Context) error {
		return deleteOne.DeleteTodo(c, store)
	})
	e.DELETE("/todos", func(c echo.Context) error {
		return deleteAll.DeleteAllTodos(c, store)
	})

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      e,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
