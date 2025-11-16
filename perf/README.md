# Производительность и оптимизация

Этот каталог содержит примеры оптимизации производительности в Go.

## Файлы

- `analysis.md` - подробная теория по производительности
- `examples.go` - примеры оптимизаций
- `benchmark_test.go` - benchmark тесты для сравнения
- `cmd/demo/main.go` - демо приложение с профилированием

## Быстрый старт

### Запуск демо приложения

```bash
go run ./perf/cmd/demo/main.go
```

Открой в браузере: http://localhost:6060/debug/pprof/

### Запуск benchmark тестов

```bash
# Все benchmark
go test ./perf -bench=.

# Конкретный benchmark
go test ./perf -bench=BenchmarkProcessDataFast

# С детальной информацией о памяти
go test ./perf -bench=. -benchmem

# С профилированием
go test ./perf -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof
```

### Профилирование

```bash
# Запусти приложение
go run ./perf/cmd/demo/main.go

# В другом терминале собери CPU профиль (30 секунд)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory профиль
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine профиль
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### Анализ профиля

В интерактивном режиме pprof:

```bash
(pprof) top          # топ функций по CPU/памяти
(pprof) list funcName # показать код функции
(pprof) web          # визуализация (нужен graphviz)
(pprof) png          # сохранить график
```

## Результаты benchmark

Из последнего запуска:

```
BenchmarkProcessDataSlow-8   739642    1676 ns/op  (неоптимизированная)
BenchmarkProcessDataFast-8   825008    1499 ns/op  (оптимизированная) ✅ быстрее!

BenchmarkBuildStringSlow-8   168463    7344 ns/op  (конкатенация)
BenchmarkBuildStringFast-8  2217252     635 ns/op  (strings.Builder) ✅ в 11 раз быстрее!

BenchmarkWithPool-8          873586    1390 ns/op  (с sync.Pool)
BenchmarkWithoutPool-8       833958    1411 ns/op  (без пула)
```

## Ключевые оптимизации

1. **Предварительное выделение памяти** - для slices и maps
2. **sync.Pool** - переиспользование объектов
3. **strings.Builder** - вместо конкатенации
4. **Кэширование** - для тяжёлых вычислений
5. **Профилирование** - поиск узких мест

## Дополнительная информация

Смотри `analysis.md` для подробной теории и best practices.
