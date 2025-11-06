package main

import (
	"day2/models"
	"day2/tools"
	"fmt"
	"day2/data"
	"time"
)

func main() {
	// Тестовые данные
	numbers := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // массив чисел

	fmt.Println("=== Статистика чисел ===")     // вывод заголовка
	fmt.Printf("Исходные числа: %v\n", numbers) // вывод исходных чисел
	fmt.Println()                               // вывод пустой строки

	// Вычисляем статистику
	stats := CalculateStats(numbers)

	// Выводим результат
	PrintStats(stats) // вывод результата

	res, err := Divide(10, 2)
	fmt.Println("Divide 10/2:", res, err)

	_, err = Divide(10, 0)
	fmt.Println("Divide 10/0:", err)

	m, err := Median([]float64{1, 3, 2, 4})
	fmt.Println("Median:", m, err)

	s := TimedSum([]float64{1, 2, 3, 4, 5})
	fmt.Println("TimedSum:", s)

	WithCleanup()
	fmt.Println("Конец программы")

	// Создаём пользователя
	user, err := models.NewUser("Алексей", "alex@example.com", "12345678", 25)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Тестируем методы
	fmt.Println(user.Introduce())
	fmt.Println("Взрослый:", user.IsAdult())
	fmt.Println(user.GetInfo())

	// Активируем
	user.Activate()
	fmt.Println("Активен:", user.IsActive)

	// Тестируем CSV→JSON конвертер
	fmt.Println("\n=== CSV → JSON конвертер ===")
	err = tools.Csv2Json("test.csv", "output.json")
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
	} else {
		fmt.Println("✅ Конвертация успешна! Проверь файл output.json")
	}

	fmt.Println("\n=== Работа с JSON ===")
	// Загружаем пользователей
	users, err := data.LoadUsers("data/users.json")
	if err != nil {
		fmt.Println("Ошибка загрузки:", err)
		// Создаём начальных пользователей
		users = []data.User{}
	}

	// Добавляем нового пользователя
	newUser := data.User{
		ID:        1,
		Name:      "Alice",
		Email:     "alice@example.com",
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	users = data.AddUser(users, newUser)

	// Сохраняем обратно
	err = data.SaveUsers("data/users.json", users)
	if err != nil {
		fmt.Println("Ошибка сохранения:", err)
	} else {
		fmt.Println("✅ Пользователи сохранены!")
	}
}
