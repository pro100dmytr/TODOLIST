package main

import (
	"TODO_List/menu"
	"bufio"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
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

var clear map[string]func()
var db *sql.DB

func main() {
	connectDB()

	e := echo.New()

	e.GET("/todos", getTodos)
	e.POST("/todos", createTodo)
	e.PUT("/todos/:id", updateTodo)
	e.DELETE("/todos/:id", deleteTodo)
	e.DELETE("/todos", deleteAllTodos)

	go func() {
		if err := e.Start(":1010"); err != nil {
			log.Fatal(err)
		}
	}()

	handleMenu()
}

func connectDB() {
	connStr := "host=localhost port=4050 user=myuser dbname=todolistdb password=4050 sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Успішно підключено до бази даних!")
}

func handleMenu() {
	userInput := bufio.NewReader(os.Stdin)

	for {
		menu.PrintMenuOptions()

		action, err := userInput.ReadString('\n')
		if err != nil {
			fmt.Println("Помилка при зчитуванні введення:", err)
			continue
		}

		action = strings.TrimSpace(action)

		switch action {
		case "1":
			createTask(userInput)
		case "2":
			listTasks()
		case "3":
			updateTask(userInput)
		case "4":
			deleteTask(userInput)
		case "5":
			deleteAllTasks()
		case "6":
			fmt.Println("Вихід з програми.")
			return
		default:
			fmt.Println("Невірна дія.")
		}
	}
}

func createTask(userInput *bufio.Reader) {
	fmt.Print("Введіть ваше завдання: ")
	title, err := userInput.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	title = strings.TrimSpace(title)

	query := "INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING task_id"
	var taskID int
	err = db.QueryRow(query, title, false).Scan(&taskID)
	if err != nil {
		fmt.Println("Помилка створення завдання:", err)
		return
	}

	fmt.Printf("Створено нове завдання з ID %d\n", taskID)
}

func listTasks() {
	rows, err := db.Query("SELECT task_id, title, completed FROM tasks")
	if err != nil {
		fmt.Println("Помилка при отриманні списку завдань:", err)
		return
	}
	defer rows.Close()

	tasks := make([]Todo, 0)
	for rows.Next() {
		var task Todo
		err := rows.Scan(&task.ID, &task.Title, &task.Completed)
		if err != nil {
			fmt.Println("Помилка при скануванні рядка:", err)
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Помилка при ітерації по результатам запиту:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("Немає завдань.\n")
		return
	}

	for _, task := range tasks {
		status := "не виконано"
		if task.Completed {
			status = "виконано"
		}
		fmt.Printf("ID: %d, Завдання: %s, Статус: %s\n", task.ID, task.Title, status)
	}
}

func updateTask(userInput *bufio.Reader) {

	fmt.Println("Усі доступні завдання\n")
	listTasks()

	fmt.Print("\nВведіть ID завдання для оновлення: ")
	idStr, err := userInput.ReadString('\n')
	if err != nil {
		fmt.Println("Помилка при зчитуванні ID завдання:", err)
		return
	}
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Невірний формат ID.")
		return
	}

	fmt.Print("Введіть нову назву завдання: ")
	newTitle, err := userInput.ReadString('\n')
	if err != nil {
		fmt.Println("Помилка при зчитуванні нової назви завдання:", err)
		return
	}
	newTitle = strings.TrimSpace(newTitle)

	fmt.Print("Чи виконано завдання? (так/ні): ")
	status, err := userInput.ReadString('\n')
	if err != nil {
		fmt.Println("Помилка при зчитуванні статусу завдання:", err)
		return
	}
	status = strings.TrimSpace(status)
	completed := (status == "так")

	query := "UPDATE tasks SET title = $1, completed = $2 WHERE task_id = $3"
	result, err := db.Exec(query, newTitle, completed, id)
	if err != nil {
		fmt.Println("Помилка при оновленні завдання:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Помилка при отриманні кількості змінених рядків:", err)
		return
	}

	if rowsAffected == 0 {
		fmt.Println("Завдання з вказаним ID не знайдено.")
	} else {
		fmt.Println("Завдання успішно оновлено.\n")
	}
}

func deleteTask(userInput *bufio.Reader) {
	fmt.Print("Введіть ID завдання для видалення: ")
	idStr, err := userInput.ReadString('\n')
	if err != nil {
		fmt.Println("Помилка при зчитуванні ID завдання:", err)
		return
	}
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Невірний формат ID.")
		return
	}

	query := "DELETE FROM tasks WHERE task_id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		fmt.Println("Помилка при видаленні завдання:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Помилка при отриманні кількості змінених рядків:", err)
		return
	}

	if rowsAffected == 0 {
		fmt.Println("Завдання з вказаним ID не знайдено.")
	} else {
		fmt.Println("Завдання успішно видалено.\n")
	}
}

func deleteAllTasks() {
	query := "DELETE FROM tasks"
	result, err := db.Exec(query)
	if err != nil {
		fmt.Println("Помилка при видаленні всіх завдань:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Помилка при отриманні кількості змінених рядків:", err)
		return
	}

	fmt.Printf("Всі завдання видалено. Кількість змінених рядків: %d\n", rowsAffected)
}

func getTodos(c echo.Context) error {
	rows, err := db.Query("SELECT task_id, title, completed FROM tasks")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при отриманні завдань"})
	}
	defer rows.Close()

	tasks := make([]Todo, 0)
	for rows.Next() {
		var task Todo
		err := rows.Scan(&task.ID, &task.Title, &task.Completed)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при скануванні рядків"})
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при ітерації по результатам запиту"})
	}

	return c.JSON(http.StatusOK, tasks)
}

func createTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Помилка при обробці запиту"})
	}

	query := "INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING task_id"
	var id int
	err := db.QueryRow(query, todo.Title, todo.Completed).Scan(&id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при додаванні завдання в базу даних"})
	}

	todo.ID = id
	return c.JSON(http.StatusCreated, todo)
}

func updateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Невірний ID"})
	}

	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Помилка при обробці запиту"})
	}

	query := "UPDATE tasks SET title = $1, completed = $2 WHERE task_id = $3"
	result, err := db.Exec(query, todo.Title, todo.Completed, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при оновленні завдання"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при отриманні кількості змінених рядків"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, todo)
}

func deleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Невірний ID"})
	}

	query := "DELETE FROM tasks WHERE task_id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при видаленні завдання"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при отриманні кількості змінених рядків"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func deleteAllTodos(c echo.Context) error {
	query := "DELETE FROM tasks"
	result, err := db.Exec(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при видаленні всіх завдань"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Помилка при отриманні кількості змінених рядків"})
	}

	if rowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}
