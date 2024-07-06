package menu

import "fmt"

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
