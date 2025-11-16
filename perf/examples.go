package perf

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // автоматически добавляет /debug/pprof
	"runtime"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// ТЕОРИЯ: Примеры оптимизации производительности
// ============================================================================

// ============================================================================
// 1. ПРОФИЛИРОВАНИЕ
// ============================================================================

// StartProfilingServer запускает HTTP сервер для профилирования
func StartProfilingServer(port string) {
	// ТЕОРИЯ: pprof автоматически добавляет endpoints:
	// - /debug/pprof/ - список всех профилей
	// - /debug/pprof/profile - CPU профиль
	// - /debug/pprof/heap - Memory профиль
	// - /debug/pprof/goroutine - Goroutine профиль
	// - /debug/pprof/block - Block профиль

	go func() {
		fmt.Printf("Профилирование доступно на http://localhost:%s/debug/pprof/\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Printf("Ошибка запуска сервера профилирования: %v\n", err)
		}
	}()
}

// PrintMemoryStats выводит статистику памяти
func PrintMemoryStats() {
	// ТЕОРИЯ: runtime.MemStats содержит информацию о памяти
	// - Alloc - текущая выделенная память
	// - TotalAlloc - всего выделено за время работы
	// - Sys - память запрошенная у ОС
	// - NumGC - количество сборок мусора

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("=== Memory Stats ===\n")
	fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
	fmt.Printf("TotalAlloc: %d KB\n", m.TotalAlloc/1024)
	fmt.Printf("Sys: %d KB\n", m.Sys/1024)
	fmt.Printf("NumGC: %d\n", m.NumGC)
	fmt.Printf("GC Pause: %v\n", time.Duration(m.PauseTotalNs))
	fmt.Printf("===================\n")
}

// ============================================================================
// 2. ОПТИМИЗАЦИЯ АЛЛОКАЦИЙ
// ============================================================================

// ProcessDataSlow НЕОПТИМИЗИРОВАННАЯ версия
func ProcessDataSlow(data []int) []int {
	// ТЕОРИЯ: Проблемы:
	// - Создаётся новый slice при каждой итерации
	// - Много аллокаций
	// - Нагрузка на GC

	result := make([]int, 0)
	for _, v := range data {
		if v%2 == 0 {
			result = append(result, v*2) // может быть переаллокация
		}
	}
	return result
}

// ProcessDataFast ОПТИМИЗИРОВАННАЯ версия
func ProcessDataFast(data []int) []int {
	// ТЕОРИЯ: Оптимизации:
	// - Предварительное выделение памяти (capacity = len(data))
	// - Меньше переаллокаций
	// - Меньше работы для GC

	result := make([]int, 0, len(data)) // предварительное выделение
	for _, v := range data {
		if v%2 == 0 {
			result = append(result, v*2) // без переаллокаций
		}
	}
	return result
}

// ============================================================================
// 3. SYNC.POOL ДЛЯ ПЕРЕИСПОЛЬЗОВАНИЯ
// ============================================================================

// BufferPool - пул буферов для переиспользования
var BufferPool = sync.Pool{
	// ТЕОРИЯ: New функция вызывается когда пул пуст
	// - Создаёт новый объект
	New: func() interface{} {
		return make([]byte, 0, 1024) // буфер 1KB
	},
}

// GetBuffer получает буфер из пула
func GetBuffer() []byte {
	// ТЕОРИЯ: Get() возвращает объект из пула
	// - Если пул пуст - вызывается New()
	// - Объект нужно привести к нужному типу
	return BufferPool.Get().([]byte)
}

// PutBuffer возвращает буфер в пул
func PutBuffer(buf []byte) {
	// ТЕОРИЯ: Put() возвращает объект в пул
	// - Важно очистить объект перед возвратом
	// - Иначе данные могут "протечь" в следующий раз

	buf = buf[:0] // очищаем, но сохраняем capacity
	BufferPool.Put(buf)
}

// ProcessWithPool использует sync.Pool для переиспользования буферов
func ProcessWithPool(data []string) []string {
	// ТЕОРИЯ: Использование sync.Pool:
	// - Уменьшает аллокации
	// - Уменьшает нагрузку на GC
	// - Особенно полезно для часто создаваемых объектов

	buf := GetBuffer()
	defer PutBuffer(buf) // возвращаем в пул

	for _, s := range data {
		buf = append(buf, []byte(s)...)
		buf = append(buf, '\n')
	}

	return strings.Split(string(buf), "\n")
}

// ============================================================================
// 4. STRINGS BUILDER VS КОНКАТЕНАЦИЯ
// ============================================================================

// BuildStringSlow НЕОПТИМИЗИРОВАННАЯ версия (конкатенация)
func BuildStringSlow(parts []string) string {
	// ТЕОРИЯ: Проблемы конкатенации:
	// - Каждая операция + создаёт новый string
	// - Много аллокаций и копирований
	// - O(n²) сложность

	var result string
	for _, part := range parts {
		result += part // создаёт новый string каждый раз!
	}
	return result
}

// BuildStringFast ОПТИМИЗИРОВАННАЯ версия (strings.Builder)
func BuildStringFast(parts []string) string {
	// ТЕОРИЯ: strings.Builder:
	// - Эффективное построение строк
	// - Минимум аллокаций
	// - O(n) сложность

	var builder strings.Builder
	builder.Grow(len(parts) * 10) // предварительное выделение (примерно)

	for _, part := range parts {
		builder.WriteString(part) // эффективно
	}

	return builder.String()
}

// ============================================================================
// 5. КЭШИРОВАНИЕ
// ============================================================================

// Cache - простой кэш с мьютексом
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewCache создаёт новый кэш
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Get получает значение из кэша
func (c *Cache) Get(key string) (interface{}, bool) {
	// ТЕОРИЯ: RLock для чтения
	// - Позволяет множественные одновременные чтения
	// - Блокирует только запись

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}

// Set устанавливает значение в кэш
func (c *Cache) Set(key string, value interface{}) {
	// ТЕОРИЯ: Lock для записи
	// - Блокирует и чтение и запись
	// - Гарантирует безопасность

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// ExpensiveOperationWithCache выполняет тяжёлую операцию с кэшированием
func ExpensiveOperationWithCache(cache *Cache, input string) int {
	// ТЕОРИЯ: Кэширование результатов:
	// - Избегаем повторных вычислений
	// - Ускоряем работу программы
	// - Особенно полезно для тяжёлых операций

	// Проверяем кэш
	if result, ok := cache.Get(input); ok {
		return result.(int) // кэш попадание - быстро!
	}

	// Тяжёлое вычисление (симуляция)
	time.Sleep(10 * time.Millisecond)
	result := len(input) * 1000

	// Сохраняем в кэш
	cache.Set(input, result)

	return result
}

// ============================================================================
// 6. НАСТРОЙКА GC
// ============================================================================

// ConfigureGC настраивает параметры Garbage Collector
func ConfigureGC() {
	// ТЕОРИЯ: Настройка GC через runtime/debug
	// - SetGCPercent - процент CPU для GC (по умолчанию 100)
	// - Больше значение = меньше паузы GC, но больше CPU
	// - Меньше значение = больше паузы GC, но меньше CPU

	// Увеличиваем процент CPU для GC (меньше паузы)
	// runtime/debug.SetGCPercent(200) // 200% = больше CPU для GC

	// Уменьшаем процент CPU для GC (больше паузы, но меньше CPU)
	// runtime/debug.SetGCPercent(50) // 50% = меньше CPU для GC

	fmt.Println("GC настроен (используй runtime/debug.SetGCPercent для изменения)")
}

// ForceGC принудительно запускает сборку мусора
func ForceGC() {
	// ТЕОРИЯ: runtime.GC() принудительно запускает GC
	// - Полезно для тестирования
	// - Не рекомендуется в продакшене (GC сам знает когда запускаться)

	runtime.GC()
	fmt.Println("GC принудительно запущен")
}

// ============================================================================
// 7. ПРИМЕР ИСПОЛЬЗОВАНИЯ
// ============================================================================

// ExampleUsage демонстрирует использование всех оптимизаций
func ExampleUsage() {
	// Запускаем сервер профилирования
	StartProfilingServer("6060")

	// Настраиваем GC
	ConfigureGC()

	// Печатаем статистику памяти
	PrintMemoryStats()

	// Пример использования sync.Pool
	data := []string{"hello", "world", "golang"}
	result := ProcessWithPool(data)
	fmt.Printf("Результат: %v\n", result)

	// Пример использования кэша
	cache := NewCache()

	start := time.Now()
	ExpensiveOperationWithCache(cache, "test1")
	fmt.Printf("Первое вычисление: %v\n", time.Since(start))

	start = time.Now()
	ExpensiveOperationWithCache(cache, "test1") // из кэша!
	fmt.Printf("Второе вычисление (из кэша): %v\n", time.Since(start))

	// Печатаем финальную статистику
	PrintMemoryStats()
}
