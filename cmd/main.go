package main

import (
	"TODO_List/internal/config"
	"TODO_List/internal/http-server/handlers/create"
	"TODO_List/internal/http-server/handlers/createcategory"
	"TODO_List/internal/http-server/handlers/deleteall"
	"TODO_List/internal/http-server/handlers/deletecategory"
	"TODO_List/internal/http-server/handlers/deleteone"
	"TODO_List/internal/http-server/handlers/get"
	"TODO_List/internal/http-server/handlers/getcategory"
	"TODO_List/internal/http-server/handlers/getcategorybyid"
	"TODO_List/internal/http-server/handlers/update"
	"TODO_List/internal/http-server/handlers/updatecategory"
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
		return deleteone.DeleteTodo(c, store)
	})
	e.DELETE("/todos", func(c echo.Context) error {
		return deleteall.DeleteAllTodos(c, store)
	})

	//categories
	e.GET("/categories", func(c echo.Context) error {
		return getcategory.GetAllCategories(c, store)
	})
	e.GET("/categories/:id", func(c echo.Context) error {
		return getcategorybyid.GetCategoryById(c, store)
	})
	e.POST("/categories", func(c echo.Context) error {
		return createcategory.CreateCategory(c, store)
	})
	e.PUT("/categories/:id", func(c echo.Context) error {
		return updatecategory.UpdateCategory(c, store)
	})
	e.DELETE("/categories/:id", func(c echo.Context) error {
		return deletecategory.DeleteCategory(c, store)
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
