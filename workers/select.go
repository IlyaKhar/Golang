package workers

import (
    "fmt"
    "time"
)

// ProcessWithTimeout обрабатывает задачу с таймаутом
func ProcessWithTimeout(task Task, timeout time.Duration) (Result, error) {
    resultCh := make(chan Result, 1)
    
    go func() {
        // Выполняем задачу
        result := ProcessTask(task)
        resultCh <- result
    }()
    
    select {
    case result := <-resultCh:
        return result, nil
    case <-time.After(timeout):
        return Result{}, fmt.Errorf("task %d timeout", task.ID)
    }
}

// MonitorWorkers мониторит несколько каналов результатов
func MonitorWorkers(results []<-chan Result, done chan bool) {
    for {
        select {
        case r := <-results[0]:
            fmt.Println("Worker 1:", r)
        case r := <-results[1]:
            fmt.Println("Worker 2:", r)
        case <-done:
            return
        }
    }
}