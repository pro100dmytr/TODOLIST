package updateTask

import (
	"TODO_List/internal/api/listTasks"
	"TODO_List/internal/storage/postgresql"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func UpdateTask(userInput *bufio.Reader) {

	fmt.Println("Усі доступні завдання\n")
	listTasks.ListTasks()

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
	result, err := postgresql.Db.Exec(query, newTitle, completed, id)
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
