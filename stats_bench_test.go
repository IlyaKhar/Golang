package main

import (
	"math/rand"
	"testing"
	"strconv"
)

func genFloats(n int) []float64 {
	a := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Float64()
	}
	return a
}

func BenchmarkCalculateStats(b *testing.B) {
	sizes := []int{1e3, 1e4, 5e4}
	for _, n := range sizes {
		b.Run(
			// имя сабтеста
			funcName(n),
			func(b *testing.B) {
				data := genFloats(n)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = CalculateStats(data)
				}
			},
		)
	}
}

func funcName(n int) string { return "N=" + strconv.Itoa(n) }