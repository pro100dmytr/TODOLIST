package user_interface

import (
	"TODO_List/internal/api/createTask"
	"TODO_List/internal/api/deleteTasks/deleteAllTasks"
	"TODO_List/internal/api/deleteTasks/deleteOneTask"
	"TODO_List/internal/api/listTasks"
	"TODO_List/internal/api/updateTask"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func HandleMenu() {
	userInput := bufio.NewReader(os.Stdin)

	for {
		PrintMenuOptions()

		action, err := userInput.ReadString('\n')
		if err != nil {
			fmt.Println("Помилка при зчитуванні введення:", err)
			continue
		}

		action = strings.TrimSpace(action)

		switch action {
		case "1":
			createTask.CreateTask(userInput)
		case "2":
			listTasks.ListTasks()
		case "3":
			updateTask.UpdateTask(userInput)
		case "4":
			deleteOneTask.DeleteTask(userInput)
		case "5":
			deleteAllTasks.DeleteAllTasks()
		case "6":
			fmt.Println("Вихід з програми.")
			return
		default:
			fmt.Println("Невірна дія.")
		}
	}
}

func PrintMenuOptions() {
	fmt.Println("Оберіть опцію:")
	fmt.Println("1. Створити нове завдання")
	fmt.Println("2. Показати список завдань")
	fmt.Println("3. Оновити завдання")
	fmt.Println("4. Видалити завдання")
	fmt.Println("5. Видалити всі завдання")
	fmt.Println("6. Вихід")
	fmt.Print("Введіть номер опції: ")
}
