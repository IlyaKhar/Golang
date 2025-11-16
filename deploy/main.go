package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GetPort читает порт из переменной окружения
func GetPort() string {
	// ТЕОРИЯ: Платформы (Render, Heroku) назначают порт через переменную PORT
	// - Если PORT не установлен (локально) - используй дефолтный порт (8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // дефолтный порт для локальной разработки
	}
	return port
}

func main() {
	// ТЕОРИЯ: Читаем порт из переменной окружения
	port := GetPort()
	log.Printf("Запуск сервера на порту %s", port)

	// ТЕОРИЯ: Настраиваем роутер
	mux := http.NewServeMux()
	
	// Health check endpoint (важен для платформ)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Главная страница
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from deployed Go app!\n")
		fmt.Fprintf(w, "Port: %s\n", port)
		fmt.Fprintf(w, "Environment: %s\n", os.Getenv("ENVIRONMENT"))
	})

	// ТЕОРИЯ: Создаём сервер с правильным адресом
	// - ":" + port означает слушать на всех интерфейсах
	// - Это нужно для деплоя (платформа сама направляет трафик)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// ТЕОРИЯ: Graceful shutdown - корректное завершение
	// - При получении сигнала (Ctrl+C, SIGTERM) сервер корректно завершится
	// - Это важно для деплоя (платформы отправляют SIGTERM при остановке)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	log.Printf("Сервер запущен на http://localhost:%s", port)

	// ТЕОРИЯ: Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// ТЕОРИЯ: Даём время завершить активные запросы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер завершён")
}

