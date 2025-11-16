package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vercel/go-bridge/go/bridge"
)

// ТЕОРИЯ: Vercel работает через serverless functions
// - Каждый HTTP запрос обрабатывается отдельной функцией
// - Нужно использовать специальный handler для Vercel
// - Это отличается от обычного HTTP сервера

func handler(w http.ResponseWriter, r *http.Request) {
	// ТЕОРИЯ: Обработчик для Vercel
	// - Vercel передаёт запросы в эту функцию
	// - Нужно обрабатывать разные пути вручную
	
	path := r.URL.Path
	
	switch path {
	case "/health":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	case "/":
		fmt.Fprintf(w, "Hello from Vercel!\n")
		fmt.Fprintf(w, "Environment: %s\n", os.Getenv("VERCEL_ENV"))
		fmt.Fprintf(w, "Path: %s\n", path)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	// ТЕОРИЯ: Vercel использует bridge для запуска Go функций
	// - bridge.ListenAndServe запускает serverless функцию
	// - Это не обычный HTTP сервер!
	bridge.Start(handler)
}

