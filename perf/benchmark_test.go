package perf

import (
	"strings"
	"testing"
)

// ============================================================================
// BENCHMARK ТЕСТЫ ДЛЯ СРАВНЕНИЯ ПРОИЗВОДИТЕЛЬНОСТИ
// ============================================================================

// BenchmarkProcessDataSlow тестирует неоптимизированную версию
func BenchmarkProcessDataSlow(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProcessDataSlow(data)
	}
}

// BenchmarkProcessDataFast тестирует оптимизированную версию
func BenchmarkProcessDataFast(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProcessDataFast(data)
	}
}

// BenchmarkBuildStringSlow тестирует конкатенацию строк
func BenchmarkBuildStringSlow(b *testing.B) {
	parts := make([]string, 100)
	for i := range parts {
		parts[i] = "test string "
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildStringSlow(parts)
	}
}

// BenchmarkBuildStringFast тестирует strings.Builder
func BenchmarkBuildStringFast(b *testing.B) {
	parts := make([]string, 100)
	for i := range parts {
		parts[i] = "test string "
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildStringFast(parts)
	}
}

// BenchmarkWithPool тестирует использование sync.Pool
func BenchmarkWithPool(b *testing.B) {
	data := make([]string, 100)
	for i := range data {
		data[i] = "test data"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ProcessWithPool(data)
	}
}

// BenchmarkWithoutPool тестирует без sync.Pool (для сравнения)
func BenchmarkWithoutPool(b *testing.B) {
	data := make([]string, 100)
	for i := range data {
		data[i] = "test data"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Создаём новый буфер каждый раз
		buf := make([]byte, 0, 1024)
		for _, s := range data {
			buf = append(buf, []byte(s)...)
			buf = append(buf, '\n')
		}
		_ = strings.Split(string(buf), "\n")
	}
}

// BenchmarkCacheHit тестирует кэш попадание
func BenchmarkCacheHit(b *testing.B) {
	cache := NewCache()
	cache.Set("test", 12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("test")
	}
}

// BenchmarkCacheMiss тестирует кэш промах
func BenchmarkCacheMiss(b *testing.B) {
	cache := NewCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("nonexistent")
	}
}
