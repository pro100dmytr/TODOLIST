package main

import (
	"TODO_List/menu"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

const filename = "todos.json"

var todos []Todo
var nextID = 1

func saveTodos() error {
	data, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func loadTodos() error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Файл не існує, починаємо з пустого списку
		}
		return err
	}

	return json.Unmarshal(data, &todos)
}

func main() {
	e := echo.New()

	if err := loadTodos(); err != nil {
		log.Fatal(err)
	}

	e.GET("/todos", lookTodos)
	e.POST("/todos", createTodo)
	e.PUT("/todos/:id", updateTodo)
	e.DELETE("/todos/:id", deleteTodo)
	e.DELETE("/todos", deleteAllTodo)

	go func() {
		e.Logger.Fatal(e.Start(":1010"))
	}()

	user := bufio.NewReader(os.Stdin)

	for {
		menu.ListPossibilities()

		action, _ := user.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "1":
			fmt.Print("Enter a name for the new task: ")

			title, err := user.ReadString('\n')

			if err != nil {
				fmt.Println("Error reading input: ", err)
				return
			}

			title = strings.TrimSpace(title)

			if title == "" {
				fmt.Println("Task name cannot be empty")
				return
			}

			newTodo := Todo{ID: nextID, Title: title, Completed: false}
			todos = append(todos, newTodo)

			fmt.Printf("Added new task\nid: %d\ntitle: %s\ncompleted: %v\n", newTodo.ID, newTodo.Title, newTodo.Completed)
			nextID++

		case "2":
			fmt.Println("list of tasks")

			if len(todos) == 0 {
				fmt.Println("No tasks found")
			}

			for _, todo := range todos {
				status := "not done"
				if todo.Completed {
					status = "done"
				}

				fmt.Printf("id: %d\ntitle: %s\ncompleted: %s \n\n", todo.ID, todo.Title, status)
			}

		case "3":
			fmt.Print("Enter an id of the task to upgrade: ")

			idStr, _ := user.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, _ := strconv.Atoi(idStr)

			for i, todo := range todos {
				if todo.ID == id {
					fmt.Print("Enter a new title of the task: ")

					newTitle, _ := user.ReadString('\n')
					newTitle = strings.TrimSpace(newTitle)
					todos[i].Title = newTitle

					fmt.Print("Task completed?")

					doneStr, _ := user.ReadString('\n')
					doneStr = strings.TrimSpace(doneStr)
					todos[i].Completed = doneStr == "done"

					fmt.Println("Update task")
					break
				}
			}

		case "4":
			fmt.Print("Enter an id of the task to delete: ")

			idStr, _ := user.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, _ := strconv.Atoi(idStr)

			for i, todo := range todos {
				status := "not done"
				if todo.Completed {
					status = "done"
				}

				if todo.ID == id {
					todos = append(todos[:i], todos[i+1:]...)

					fmt.Println("\nDelete task\n")
					fmt.Printf("id: %d\ntitle: %s\ncompleted: %s\n\n", todo.ID, todo.Title, status)
					break
				}
			}

		case "5":
			fmt.Println("Deleting all tasks")

			todos = []Todo{}
			nextID = 1
			saveTodos()

			fmt.Println("All tasks deleted")

		case "6":
			defer saveTodos()
			fmt.Println("You have a todo list")
			return

		default:
			fmt.Println("Invalid action")
		}
	}

}

func lookTodos(c echo.Context) error {
	return c.JSON(http.StatusOK, todos)
}

func createTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(&todo); err != nil {
		return err
	}

	todo.ID = nextID
	nextID++
	todos = append(todos, *todo)

	return c.JSON(http.StatusCreated, todo)
}

func updateTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Title = c.FormValue("title")
			todos[i].Completed, _ = strconv.ParseBool(c.FormValue("completed"))
			return c.JSON(http.StatusOK, todos[i])
		}
	}

	return c.NoContent(http.StatusNotFound)
}

func deleteTodo(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	for i, todo := range todos {

		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}

	return c.NoContent(http.StatusNotFound)
}

func deleteAllTodo(c echo.Context) error {

	todos = []Todo{}
	nextID = 1
	if err := saveTodos(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save todos"})
	}

	return c.NoContent(http.StatusNoContent)
}
