package main

import (
	"math"
	"testing"
)

// TestCalculateStats тестирует функцию CalculateStats
func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []float64
		expected Stats
	}{
		{
			name:     "Пустой slice",
			numbers:  []float64{},
			expected: Stats{},
		},
		{
			name:    "Одно число",
			numbers: []float64{5.0},
			expected: Stats{
				Count:    1,
				Sum:      5.0,
				Mean:     5.0,
				Min:      5.0,
				Max:      5.0,
				Variance: 0.0,
				StdDev:   0.0,
			},
		},
		{
			name:    "Два числа",
			numbers: []float64{1.0, 3.0},
			expected: Stats{
				Count:    2,
				Sum:      4.0,
				Mean:     2.0,
				Min:      1.0,
				Max:      3.0,
				Variance: 1.0, // (1-2)² + (3-2)² = 1 + 1 = 2, 2/2 = 1
				StdDev:   1.0, // √1 = 1
			},
		},
		{
			name:    "Несколько чисел",
			numbers: []float64{1, 2, 3, 4, 5},
			expected: Stats{
				Count:    5,
				Sum:      15.0,
				Mean:     3.0,
				Min:      1.0,
				Max:      5.0,
				Variance: 2.0,            // (1-3)² + (2-3)² + (3-3)² + (4-3)² + (5-3)² = 4+1+0+1+4 = 10, 10/5 = 2
				StdDev:   math.Sqrt(2.0), // √2 ≈ 1.414
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStats(tt.numbers)

			// Проверяем каждое поле с учетом погрешности для float64
			if result.Count != tt.expected.Count {
				t.Errorf("Count: получили %d, ожидали %d", result.Count, tt.expected.Count)
			}

			if !almostEqual(result.Sum, tt.expected.Sum) {
				t.Errorf("Sum: получили %f, ожидали %f", result.Sum, tt.expected.Sum)
			}

			if !almostEqual(result.Mean, tt.expected.Mean) {
				t.Errorf("Mean: получили %f, ожидали %f", result.Mean, tt.expected.Mean)
			}

			if !almostEqual(result.Min, tt.expected.Min) {
				t.Errorf("Min: получили %f, ожидали %f", result.Min, tt.expected.Min)
			}

			if !almostEqual(result.Max, tt.expected.Max) {
				t.Errorf("Max: получили %f, ожидали %f", result.Max, tt.expected.Max)
			}

			if !almostEqual(result.Variance, tt.expected.Variance) {
				t.Errorf("Variance: получили %f, ожидали %f", result.Variance, tt.expected.Variance)
			}

			if !almostEqual(result.StdDev, tt.expected.StdDev) {
				t.Errorf("StdDev: получили %f, ожидали %f", result.StdDev, tt.expected.StdDev)
			}
		})
	}
}

// almostEqual сравнивает два float64 с учетом погрешности
func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

// TestCalculateStatsEdgeCases тестирует граничные случаи
func TestCalculateStatsEdgeCases(t *testing.T) {
	// Тест с отрицательными числами
	numbers := []float64{-1, -2, -3, -4, -5}
	stats := CalculateStats(numbers)

	if stats.Count != 5 {
		t.Errorf("Count с отрицательными числами: получили %d, ожидали 5", stats.Count)
	}

	if !almostEqual(stats.Sum, -15.0) {
		t.Errorf("Sum с отрицательными числами: получили %f, ожидали -15.0", stats.Sum)
	}

	if !almostEqual(stats.Mean, -3.0) {
		t.Errorf("Mean с отрицательными числами: получили %f, ожидали -3.0", stats.Mean)
	}

	if !almostEqual(stats.Min, -5.0) {
		t.Errorf("Min с отрицательными числами: получили %f, ожидали -5.0", stats.Min)
	}

	if !almostEqual(stats.Max, -1.0) {
		t.Errorf("Max с отрицательными числами: получили %f, ожидали -1.0", stats.Max)
	}
}
