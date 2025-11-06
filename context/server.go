package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

type App struct {
	mu      sync.RWMutex
	counter int
	srv     *http.Server
}

func (a *App) handleGetCounter(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	val := a.counter
	a.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"counter": val})
}

func (a *App) handleIncCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mu.Lock()
	a.counter++
	val := a.counter
	a.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"counter": val})
}

func (a *App) handleWork(w http.ResponseWriter, r *http.Request) {
	// таймаут из query, по умолчанию 1s
	timeoutMs := 1000
	if qs := r.URL.Query().Get("timeout"); qs != "" {
		if v, err := strconv.Atoi(qs); err == nil && v > 0 {
			timeoutMs = v
		}
	}
	// длительность «работы», по умолчанию 1500ms
	workMs := 1500
	if qs := r.URL.Query().Get("ms"); qs != "" {
		if v, err := strconv.Atoi(qs); err == nil && v > 0 {
			workMs = v
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	done := make(chan struct{}, 1)
	go func() {
		// имитируем долгую работу
		time.Sleep(time.Duration(workMs) * time.Millisecond)
		done <- struct{}{}
	}()

	select {
	case <-done:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	case <-ctx.Done():
		http.Error(w, fmt.Sprintf("timeout/cancel: %v", ctx.Err()), http.StatusGatewayTimeout)
	}
}

func (a *App) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/counter", a.handleGetCounter)
	mux.HandleFunc("/counter/incr", a.handleIncCounter)
	mux.HandleFunc("/work", a.handleWork)
	return mux
}

func (a *App) start(addr string) error {
	mux := a.routes()
	a.srv = &http.Server{Addr: addr, Handler: mux}
	log.Printf("server on %s", addr)
	return a.srv.ListenAndServe()
}

func (a *App) shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}

func main() {
	app := &App{}
	// запускаем сервер
	go func() {
		if err := app.start(":8081"); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	app.srv = &http.Server{Addr: ":8081", Handler: mux}

	// ловим Ctrl+C и завершаем
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
