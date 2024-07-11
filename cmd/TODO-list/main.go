package main

import (
	"TODO_List/internal/http-server/handlers/create"
	"TODO_List/internal/http-server/handlers/delete/deleteAll"
	"TODO_List/internal/http-server/handlers/delete/deleteOne"
	"TODO_List/internal/http-server/handlers/get"
	"TODO_List/internal/http-server/handlers/update"
	User_interface "TODO_List/user-interface"

	"TODO_List/internal/storage/postgresql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	postgresql.ConnectDB()

	e := echo.New()

	e.GET("/todos", get.GetTodos)
	e.POST("/todos", create.CreateTodo)
	e.PUT("/todos/:id", update.UpdateTodo)
	e.DELETE("/todos/:id", deleteOne.DeleteTodo)
	e.DELETE("/todos", deleteAll.DeleteAllTodos)

	go func() {
		if err := e.Start(":1010"); err != nil {
			log.Fatal(err)
		}
	}()

	User_interface.HandleMenu()
}
