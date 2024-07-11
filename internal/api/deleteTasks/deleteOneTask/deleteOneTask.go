package deleteOneTask

import (
	"TODO_List/internal/storage/postgresql"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func DeleteTask(userInput *bufio.Reader) {
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
	result, err := postgresql.Db.Exec(query, id)
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
