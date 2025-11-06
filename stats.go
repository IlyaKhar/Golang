package main

import (
	"fmt"
	"math"
)

// Stats структура для хранения статистики
type Stats struct {
	Count    int     // количество чисел
	Sum      float64 // сумма
	Mean     float64 // среднее арифметическое
	Min      float64 // минимальное значение
	Max      float64 // максимальное значение
	Variance float64 // дисперсия
	StdDev   float64 // стандартное отклонение
}

// CalculateStats вычисляет статистику для slice чисел
func CalculateStats(numbers []float64) Stats {
	// Проверяем на пустой slice
	if len(numbers) == 0 {
		return Stats{}
	}

	// Инициализация структуры Stats (используем : вместо :=)
	stats := Stats{
		Count: len(numbers), // количество чисел
		Sum:   0.0,          // сумма
		Mean:  0.0,          // среднее арифметическое
		Min:   numbers[0],   // минимальное значение (первое число)
		Max:   numbers[0],   // максимальное значение (первое число)
	}

	// Первый проход - находим сумму, минимум и максимум
	for _, number := range numbers { // проходим по всем числам в slice
		stats.Sum += number
		stats.Min = math.Min(stats.Min, number)
		stats.Max = math.Max(stats.Max, number)
	}

	// Вычисляем среднее
	stats.Mean = stats.Sum / float64(stats.Count)

	// Второй проход - вычисляем дисперсию
	for _, number := range numbers { // проходим по всем числам в slice
		stats.Variance += math.Pow(number-stats.Mean, 2)
	}

	// Завершаем вычисления
	stats.Variance = stats.Variance / float64(stats.Count)
	stats.StdDev = math.Sqrt(stats.Variance)

	return stats
}

// PrintStats красиво выводит статистику
func PrintStats(stats Stats) {
	// TODO: твоя реализация здесь
	fmt.Printf("Количество чисел: %d\n", stats.Count)        // вывод количества чисел
	fmt.Printf("Сумма: %f\n", stats.Sum)                     // вывод суммы
	fmt.Printf("Среднее арифметическое: %f\n", stats.Mean)   // вывод среднего арифметического
	fmt.Printf("Минимальное значение: %f\n", stats.Min)      // вывод минимального значения
	fmt.Printf("Максимальное значение: %f\n", stats.Max)     // вывод максимального значения
	fmt.Printf("Дисперсия: %f\n", stats.Variance)            // вывод дисперсии
	fmt.Printf("Стандартное отклонение: %f\n", stats.StdDev) // вывод стандартного отклонения
}
