package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// FanOut создаёт N воркеров, каждый обрабатывает задачи из input канала
// worker - функция, которая обрабатывает задачу и возвращает результат
// Результаты каждого воркера пишутся в отдельный канал (для Fan-in)
func FanOut[T any, R any](input <-chan T, numWorkers int, worker func(T) R) []<-chan R {
	// Создаём слайс каналов для результатов каждого воркера
	outputs := make([]<-chan R, numWorkers)
	outputChannels := make([]chan R, numWorkers)

	// Создаём каналы для каждого воркера
	for i := 0; i < numWorkers; i++ {
		ch := make(chan R)
		outputChannels[i] = ch
		outputs[i] = ch
	}

	// Запускаем воркеров
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		output := outputChannels[i]

		go func() {
			defer wg.Done()
			defer close(output) // Закрываем канал когда воркер завершится

			// Читаем задачи из общего input канала
			for task := range input {
				// Обрабатываем задачу и пишем результат в свой канал
				result := worker(task)
				output <- result
			}
		}()
	}

	// Закрываем каналы после завершения всех воркеров (опционально)
	// Но обычно это делается автоматически через defer close в горутине
	// Запускаем горутину которая ждёт завершения воркеров
	go func() {
		wg.Wait()
	}()

	return outputs
}

// FanIn собирает результаты из множества каналов в один канал
func FanIn[T any](channels []<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	// Запускаем горутину для каждого канала
	for _, ch := range channels {
		wg.Add(1)

		go func(input <-chan T) {
			defer wg.Done()

			// Читаем из своего канала и пишем в общий out
			for value := range input {
				out <- value
			}
		}(ch)
	}

	// Закрываем out после того как все горутины закончат
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

type Task struct {
	ID   int
	Data string
}

// Result представляет результат обработки задачи
type Result struct {
	TaskID int
	Output string
	Error  error
}

// WorkerPool - продвинутый пул воркеров с контекстом и сбором результатов
type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	results    chan Result
	wg         sync.WaitGroup
}

// NewWorkerPool создаёт новый пул воркеров
func NewWorkerPool(numWorkers int, taskQueueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Task, taskQueueSize),
		results:    make(chan Result, taskQueueSize),
	}
}

// Start запускает пул воркеров
// Каждый воркер читает задачи из p.tasks и обрабатывает их через processor
func (p *WorkerPool) Start(ctx context.Context, processor func(Task) Result) {
	// Запускаем p.numWorkers горутин-воркеров
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1) // Увеличиваем счётчик ожидаемых горутин

		go func() {
			defer p.wg.Done() // Уменьшаем счётчик когда горутина завершится

			// Бесконечный цикл: читаем задачи пока не отменён контекст или не закрыт канал
			for {
				select {
				case <-ctx.Done():
					// Контекст отменён - завершаем работу
					return
				case task, ok := <-p.tasks:
					// Читаем задачу из канала
					if !ok {
						// Канал закрыт - завершаем работу
						return
					}
					result := processor(task)
					p.results <- result
				}
			}
		}()
	}
}

// Submit отправляет задачу в пул
func (p *WorkerPool) Submit(task Task) error {
	select {
	case p.tasks <- task:
		// Задача успешно отправлена
		return nil
	default:
		// Буфер переполнен - возвращаем ошибку
		return fmt.Errorf("task queue is full")
	}
}

// Results возвращает канал результатов
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

// Stop останавливает пул воркеров (graceful shutdown)
func (p *WorkerPool) Stop() {
	// Закрываем канал задач - воркеры закончат обработку текущих задач и завершатся
	close(p.tasks)

	// Ждём завершения всех воркеров
	p.wg.Wait()

	// Закрываем канал результатов
	close(p.results)
}

// ProcessNumber обрабатывает число (симуляция работы)
func ProcessNumber(n int) int {
	time.Sleep(10 * time.Millisecond) // симуляция работы
	return n * n
}

// ExampleFanOutFanIn демонстрирует паттерн Fan-out/Fan-in
func ExampleFanOutFanIn() {
	// Создаём входной канал с числами
	input := make(chan int, 10)

	// Заполняем канал
	go func() {
		defer close(input)
		for i := 1; i <= 100; i++ {
			input <- i
		}
	}()

	// FanOut: создаём 5 воркеров, каждый обрабатывает числа через ProcessNumber
	workerChannels := FanOut(input, 5, ProcessNumber)

	// FanIn: собираем результаты из всех каналов в один
	results := FanIn(workerChannels)

	// Читаем и выводим результаты
	fmt.Println("Результаты обработки:")
	for result := range results {
		fmt.Printf("Результат: %d\n", result)
	}

	fmt.Println("Fan-out/Fan-in example completed")
}

// ProcessTask обрабатывает задачу (симуляция работы)
func ProcessTask(task Task) Result {
	time.Sleep(50 * time.Millisecond) // симуляция работы
	return Result{
		TaskID: task.ID,
		Output: fmt.Sprintf("processed: %s", task.Data),
		Error:  nil,
	}
}

// ExampleWorkerPool демонстрирует продвинутый worker pool
func ExampleWorkerPool() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := NewWorkerPool(3, 10) // 3 воркера, буфер на 10 задач

	// Запускаем пул воркеров
	pool.Start(ctx, ProcessTask)

	// Отправляем задачи в горутине
	var sendWg sync.WaitGroup
	sendWg.Add(1)
	go func() {
		defer sendWg.Done()
		for i := 1; i <= 20; i++ {
			task := Task{ID: i, Data: fmt.Sprintf("task-%d", i)}
			if err := pool.Submit(task); err != nil {
				fmt.Printf("failed to submit task %d: %v\n", i, err)
				break
			}
		}
	}()

	// Читаем результаты в горутине (чтобы не блокировать отправку задач)
	var readWg sync.WaitGroup
	readWg.Add(1)
	go func() {
		defer readWg.Done()
		fmt.Println("Результаты обработки задач:")
		resultCount := 0
		for result := range pool.Results() {
			resultCount++
			if result.Error != nil {
				fmt.Printf("Task %d error: %v\n", result.TaskID, result.Error)
			} else {
				fmt.Printf("Task %d: %s\n", result.TaskID, result.Output)
			}
		}
		fmt.Printf("Обработано задач: %d\n", resultCount)
	}()

	// Ждём завершения отправки всех задач
	sendWg.Wait()

	// Останавливаем пул (закроет канал задач, воркеры обработают оставшиеся задачи и завершатся)
	// Stop() также закроет канал результатов, что завершит цикл for range в горутине чтения
	pool.Stop()

	// Ждём завершения чтения всех результатов
	readWg.Wait()

	fmt.Println("Worker pool example completed")
}
