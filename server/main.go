package main

import (
	"day2/data"
	"day2/handlers"
	"day2/workers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Регистрируем обработчики
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/users", usersHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод - только GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "message": "Server is running"}`)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// GET /users - получить всех пользователей
		users, err := data.LoadUsers("data/users.json")
		if err != nil {
			http.Error(w, "Failed to load users", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			http.Error(w, "Failed to encode users", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		// Используем handler из пакета handlers
		handlers.PostUserHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobs := make(chan workers.Task, 10)
	results := make(chan workers.Result, 10)

	// Отправляем задачи
	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- workers.Task{ID: i, Data: fmt.Sprintf("task-%d", i)}
		}
		close(jobs)
	}()

	// Создаём пул
	workers.CreatePool(3, jobs, results)

	// Читаем результаты
	for i := 1; i <= 5; i++ {
		result := <-results
		fmt.Println(result)
	}

}

