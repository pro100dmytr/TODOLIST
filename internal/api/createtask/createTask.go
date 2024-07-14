package createtask

import (
	"TODO_List/internal/storage/postgresql"
	"bufio"
	"fmt"
	"strings"
)

func CreateTask(userInput *bufio.Reader, store *postgresql.Storage) {
	fmt.Print("Введіть ваше завдання: ")
	title, err := userInput.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	title = strings.TrimSpace(title)

	query := "INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING task_id"
	var taskID int
	err = store.DB.QueryRow(query, title, false).Scan(&taskID)
	if err != nil {
		fmt.Println("Помилка створення завдання:", err)
		return
	}

	fmt.Printf("Створено нове завдання з ID %d\n", taskID)
}
