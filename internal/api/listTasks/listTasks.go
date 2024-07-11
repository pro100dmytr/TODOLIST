package listTasks

import (
	"TODO_List/internal/storage/postgresql"
	"TODO_List/model"
	"fmt"
)

func ListTasks() {
	rows, err := postgresql.Db.Query("SELECT task_id, title, completed FROM tasks")
	if err != nil {
		fmt.Println("Помилка при отриманні списку завдань:", err)
		return
	}
	defer rows.Close()

	tasks := make([]model.Todo, 0)
	for rows.Next() {
		var task model.Todo
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
