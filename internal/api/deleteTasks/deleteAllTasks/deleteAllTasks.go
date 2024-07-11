package deleteAllTasks

import (
	"TODO_List/internal/storage/postgresql"
	"fmt"
)

func DeleteAllTasks() {
	query := "DELETE FROM tasks"
	result, err := postgresql.Db.Exec(query)
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
