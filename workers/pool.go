package workers

import (
	"fmt"
	"time"
)

type Task struct {
	ID   int
	Data string // что-то что нужно обработать
}

type Result struct {
	TaskID int
	Output string // результат обработки
	Err    error
}

func Worker(id int, jobs <-chan Task, results chan<- Result) {
	for task := range jobs {
		// Обработка задачи
		result := ProcessTask(task)
		// Отправка результата в results
		results <- result
	}
	fmt.Printf("Worker %d finished\n", id)
}

func CreatePool(numWorkers int, jobs <-chan Task, results chan<- Result) {
	for i := 1; i <= numWorkers; i++ {
		go Worker(i, jobs, results)
	}
}

func ProcessTask(task Task) Result {
	fmt.Printf("Processing task %d: %s\n", task.ID, task.Data)
	time.Sleep(time.Second * 2)
	fmt.Printf("Task %d processed\n", task.ID)
	return Result{
		TaskID: task.ID,
		Output: task.Data,
		Err:    nil,
	}
}
