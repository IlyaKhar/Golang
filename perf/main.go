package main

import (
	"day2/perf"
	"time"
)

func main() {
	// ТЕОРИЯ: Демонстрация оптимизаций производительности
	
	// Запускаем сервер профилирования
	// Открой в браузере: http://localhost:6060/debug/pprof/
	perf.StartProfilingServer("6060")
	
	// Даём время запуститься серверу
	time.Sleep(1 * time.Second)
	
	// Запускаем примеры
	perf.ExampleUsage()
	
	// Держим программу запущенной для профилирования
	// Нажми Ctrl+C для остановки
	select {}
}

