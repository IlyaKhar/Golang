# Производительность и оптимизация в Go

## ТЕОРИЯ: Производительность в Go

### ЧТО ТАКОЕ ПРОИЗВОДИТЕЛЬНОСТЬ:
Производительность - это скорость выполнения программы и эффективность использования ресурсов (память, CPU).

### ЗАЧЕМ ОПТИМИЗИРОВАТЬ:
1. **Быстрее работает** - пользователи не ждут
2. **Меньше ресурсов** - экономия денег на серверах
3. **Масштабируемость** - приложение выдержит большую нагрузку
4. **Лучший UX** - отзывчивое приложение

### КЛЮЧЕВЫЕ ОБЛАСТИ ОПТИМИЗАЦИИ:
1. **Garbage Collector (GC)** - управление памятью
2. **Профилирование** - поиск узких мест (bottlenecks)
3. **Алгоритмы** - выбор правильных структур данных
4. **Конкурентность** - эффективное использование goroutines

---

## 1. GARBAGE COLLECTOR (GC) В GO

### ТЕОРИЯ: Что такое GC

**Garbage Collector** - это автоматический сборщик мусора, который освобождает неиспользуемую память.

**Как работает:**
1. Программа выделяет память (создаёт объекты)
2. GC отслеживает какие объекты используются
3. Периодически GC останавливает программу (STW - Stop The World)
4. Удаляет неиспользуемые объекты
5. Освобождает память

**Проблемы GC:**
- **STW паузы** - программа останавливается во время сборки мусора
- **CPU overhead** - GC тратит процессорное время
- **Память** - GC может держать память дольше чем нужно

### КАК ОПТИМИЗИРОВАТЬ GC:

#### 1. Уменьши количество аллокаций

**Плохо:**
```go
// Создаётся новый slice при каждой итерации
for i := 0; i < 1000; i++ {
    result := make([]int, 0)  // аллокация!
    // ...
}
```

**Хорошо:**
```go
// Переиспользуем slice
result := make([]int, 0, 1000)  // предварительное выделение
for i := 0; i < 1000; i++ {
    result = result[:0]  // переиспользование
    // ...
}
```

#### 2. Используй sync.Pool для переиспользования объектов

**ТЕОРИЯ:** `sync.Pool` - пул объектов для переиспользования
- Уменьшает аллокации
- Уменьшает нагрузку на GC
- Особенно полезно для часто создаваемых объектов

**Пример:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024)
    },
}

func getBuffer() []byte {
    return bufferPool.Get().([]byte)
}

func putBuffer(buf []byte) {
    buf = buf[:0]  // очищаем
    bufferPool.Put(buf)
}
```

#### 3. Избегай утечек памяти

**Проблема:**
```go
// Глобальная переменная растёт бесконечно
var cache = make(map[string]interface{})

func addToCache(key string, value interface{}) {
    cache[key] = value  // память никогда не освободится!
}
```

**Решение:**
```go
// Ограниченный кэш с TTL
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
    
    // Автоматическое удаление через TTL
    time.AfterFunc(ttl, func() {
        c.Delete(key)
    })
}
```

#### 4. Используй указатели осознанно

**ТЕОРИЯ:** 
- Маленькие структуры (< 128 байт) лучше передавать по значению
- Большие структуры лучше передавать по указателю
- Указатели создают дополнительную нагрузку на GC

**Плохо:**
```go
type SmallStruct struct {
    a, b, c int
}

// Указатель для маленькой структуры - лишняя аллокация
func process(s *SmallStruct) { }
```

**Хорошо:**
```go
// Передаём по значению - нет аллокации
func process(s SmallStruct) { }
```

#### 5. Настрой переменные окружения GC

```bash
# Увеличить процент CPU для GC (по умолчанию 25%)
export GOGC=100  # 100% = использовать столько же CPU сколько и программа

# Отключить GC (только для тестов!)
export GOGC=off
```

---

## 2. ПРОФИЛИРОВАНИЕ (PROFILING)

### ТЕОРИЯ: Что такое профилирование

**Профилирование** - это анализ работы программы для поиска узких мест (bottlenecks).

**Инструменты Go:**
1. **pprof** - встроенный профилировщик
2. **go tool pprof** - анализ профилей
3. **benchmark** - тесты производительности

### КАК ИСПОЛЬЗОВАТЬ PPROF:

#### 1. Добавь импорт

```go
import (
    _ "net/http/pprof"  // автоматически добавляет /debug/pprof
    "net/http"
)
```

#### 2. Запусти HTTP сервер для профилирования

```go
func main() {
    // Запускаем HTTP сервер для профилирования
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Твой основной код
    // ...
}
```

#### 3. Собирай профили

```bash
# CPU профиль (30 секунд)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory профиль
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine профиль
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

#### 4. Анализируй профиль

В интерактивном режиме pprof:
```bash
(pprof) top          # топ функций по CPU/памяти
(pprof) list funcName # показать код функции
(pprof) web          # визуализация (нужен graphviz)
(pprof) png          # сохранить график
```

### ПРИМЕРЫ ПРОФИЛИРОВАНИЯ:

#### CPU профилирование

```go
import (
    "net/http"
    _ "net/http/pprof"
    "runtime"
)

func main() {
    // Включаем профилирование
    runtime.SetCPUProfileRate(100)  // 100Hz
    
    // HTTP сервер для профилирования
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // Твой код
    // ...
}
```

#### Memory профилирование

```go
import (
    "net/http"
    _ "net/http/pprof"
    "runtime"
)

func main() {
    // Включаем профилирование памяти
    runtime.MemProfileRate = 1  // профилировать каждую аллокацию
    
    // HTTP сервер
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // Твой код
    // ...
}
```

---

## 3. BENCHMARK ТЕСТЫ

### ТЕОРИЯ: Benchmark тесты

**Benchmark** - это тесты производительности, которые измеряют скорость выполнения кода.

### КАК ПИСАТЬ BENCHMARK:

```go
func BenchmarkFunction(b *testing.B) {
    // Подготовка данных
    data := prepareData()
    
    // Сброс таймера (подготовка не учитывается)
    b.ResetTimer()
    
    // Запускаем функцию b.N раз
    for i := 0; i < b.N; i++ {
        functionToTest(data)
    }
}
```

### ЗАПУСК BENCHMARK:

```bash
# Запустить все benchmark
go test -bench=.

# Запустить конкретный benchmark
go test -bench=BenchmarkFunction

# С профилированием
go test -bench=BenchmarkFunction -cpuprofile=cpu.prof
go test -bench=BenchmarkFunction -memprofile=mem.prof

# Сравнить производительность
go test -bench=BenchmarkFunction -benchmem
```

### ПРИМЕР BENCHMARK:

```go
func BenchmarkAppend(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        slice := make([]int, 0)
        for j := 0; j < 1000; j++ {
            slice = append(slice, j)
        }
    }
}

func BenchmarkAppendPreallocated(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        slice := make([]int, 0, 1000)  // предварительное выделение
        for j := 0; j < 1000; j++ {
            slice = append(slice, j)
        }
    }
}
```

---

## 4. BEST PRACTICES ДЛЯ ПРОИЗВОДИТЕЛЬНОСТИ

### 1. Предварительное выделение памяти

**Плохо:**
```go
slice := make([]int, 0)
for i := 0; i < 1000; i++ {
    slice = append(slice, i)  // переаллокация при росте
}
```

**Хорошо:**
```go
slice := make([]int, 0, 1000)  // предварительное выделение
for i := 0; i < 1000; i++ {
    slice = append(slice, i)  // без переаллокаций
}
```

### 2. Избегай лишних копирований

**Плохо:**
```go
func process(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)  // копирование
    // обработка
    return result
}
```

**Хорошо:**
```go
func process(data []int) {
    // обработка на месте, без копирования
    for i := range data {
        data[i] = data[i] * 2
    }
}
```

### 3. Используй strings.Builder вместо конкатенации

**Плохо:**
```go
var result string
for i := 0; i < 1000; i++ {
    result += strconv.Itoa(i)  // создаёт новый string каждый раз
}
```

**Хорошо:**
```go
var builder strings.Builder
for i := 0; i < 1000; i++ {
    builder.WriteString(strconv.Itoa(i))  // эффективно
}
result := builder.String()
```

### 4. Кэшируй результаты тяжёлых вычислений

```go
var cache = make(map[string]int)
var mu sync.RWMutex

func expensiveCalculation(input string) int {
    mu.RLock()
    if result, ok := cache[input]; ok {
        mu.RUnlock()
        return result  // кэш попадание
    }
    mu.RUnlock()
    
    // Тяжёлое вычисление
    result := doHeavyWork(input)
    
    mu.Lock()
    cache[input] = result
    mu.Unlock()
    
    return result
}
```

### 5. Используй правильные структуры данных

- **map** - O(1) поиск, но больше памяти
- **slice** - O(n) поиск, но меньше памяти
- **Выбор зависит от размера данных и частоты операций**

---

## 5. ИНСТРУМЕНТЫ ДЛЯ АНАЛИЗА

### go tool pprof

```bash
# CPU профиль
go tool pprof http://localhost:6060/debug/pprof/profile

# Memory профиль
go tool pprof http://localhost:6060/debug/pprof/heap

# Сохранить в файл
go tool pprof -png http://localhost:6060/debug/pprof/profile > cpu.png
```

### go tool trace

```go
import (
    "runtime/trace"
    "os"
)

func main() {
    f, _ := os.Create("trace.out")
    defer f.Close()
    
    trace.Start(f)
    defer trace.Stop()
    
    // Твой код
}
```

```bash
go tool trace trace.out
```

### go test -bench

```bash
# Benchmark с детальной информацией
go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof

# Сравнить два benchmark
go test -bench=Old -benchmem > old.txt
go test -bench=New -benchmem > new.txt
benchcmp old.txt new.txt
```

---

## 6. МЕТРИКИ ДЛЯ МОНИТОРИНГА

### Runtime метрики

```go
import "runtime"

func printStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("Sys: %d KB\n", m.Sys/1024)
    fmt.Printf("NumGC: %d\n", m.NumGC)
}
```

### GC статистика

```go
import (
    "runtime"
    "runtime/debug"
)

func printGCStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("GC runs: %d\n", m.NumGC)
    fmt.Printf("GC pause: %v\n", time.Duration(m.PauseTotalNs))
    
    // Настройки GC
    debug.SetGCPercent(100)  // изменить процент CPU для GC
}
```

---

## 7. ПРАКТИЧЕСКИЕ СОВЕТЫ

### 1. Измеряй перед оптимизацией

**Правило:** Не оптимизируй то, что не измерено!
- Используй профилирование для поиска узких мест
- Оптимизируй только то, что действительно медленно

### 2. Оптимизируй горячие пути

**Горячий путь** - код, который выполняется чаще всего.
- Сфокусируйся на оптимизации горячих путей
- Не трать время на оптимизацию редко выполняемого кода

### 3. Используй правильные алгоритмы

- Выбор алгоритма важнее микрооптимизаций
- O(n log n) лучше чем O(n²) для больших данных

### 4. Тестируй на реальных данных

- Benchmark на маленьких данных может вводить в заблуждение
- Тестируй на данных, похожих на продакшен

---

## РЕЗЮМЕ

### Ключевые принципы:

1. ✅ **Измеряй перед оптимизацией** - используй профилирование
2. ✅ **Уменьшай аллокации** - меньше работы для GC
3. ✅ **Предварительное выделение** - для slices и maps
4. ✅ **Переиспользование объектов** - sync.Pool
5. ✅ **Правильные структуры данных** - выбор важен
6. ✅ **Кэширование** - для тяжёлых вычислений
7. ✅ **Мониторинг** - отслеживай метрики в продакшене

### Инструменты:

- `pprof` - профилирование CPU и памяти
- `go test -bench` - benchmark тесты
- `go tool trace` - анализ выполнения
- `runtime.ReadMemStats` - метрики памяти

**Profit: быстрое и эффективное приложение! 🚀**
