package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

// Divide делит a на b и возвращает ошибку, если b == 0.
// Жизненный пример: нельзя делить счёт пополам, если друзей ноль.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("деление на ноль запрещено")
	}
	return a / b, nil
}

// Sqrt возвращает квадратный корень из value.
// Ошибка, если передали отрицательное число.
// Жизненный пример: корень из долга — не существует в реальном мире.
func Sqrt(value float64) (float64, error) {
	if value < 0 {
		return 0, errors.New("квадратный корень из отрицательного числа не определён")
	}
	return math.Sqrt(value), nil
}

// Median считает медиану по срезу чисел.
// Используем копию, чтобы не мутировать исходный срез.
func Median(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("невозможно найти медиану пустого набора")
	}
	copySlice := make([]float64, len(numbers))
	copy(copySlice, numbers)
	sort.Float64s(copySlice)

	mid := len(copySlice) / 2
	if len(copySlice)%2 == 1 {
		return copySlice[mid], nil
	}
	return (copySlice[mid-1] + copySlice[mid]) / 2, nil
}

// TimedSum суммирует числа и показывает, сколько заняла операция.
// Здесь defer — как «не забудь выключить свет, когда выйдешь».
func TimedSum(numbers []float64) (sum float64) {
	start := time.Now()
	// defer откладывает вызов функции до выхода из текущей функции
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("Сумма посчитана за %v\n", elapsed)
	}()

	for _, n := range numbers {
		sum += n
	}
	return
}

// WithCleanup демонстрирует типичный паттерн defer для освобождения ресурса.
// Здесь мы эмулируем «открытие ресурса» и гарантируем его «закрытие» через defer.
func WithCleanup() {
	fmt.Println("Открываем ресурс...")
	// Гарантируем, что «закроем ресурс» при выходе из функции
	defer fmt.Println("Закрываем ресурс (defer сработал)")

	fmt.Println("Работаем с ресурсом...")
	// Любые return/panic не помешают выполнению defer выше
}
