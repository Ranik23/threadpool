package threadpool

import (
	"fmt"
	"sync"
	"threadpool/internal/task"
	"time"
)

var (
	ErrQueueFull = fmt.Errorf("queue is full")
)


type ThreadPool struct {
	queueSize	int
	numWorkers  int
	taskQueue 	chan task.Task
	WorkerPool  chan chan task.Task
	CloseHandle chan bool
	mutex 		sync.Mutex
}

func NewThreadPool(numWorkers, queueSize int) *ThreadPool {
	return &ThreadPool{
		queueSize: queueSize,
		numWorkers: numWorkers,
		taskQueue: make(chan task.Task, queueSize),
		WorkerPool: make(chan chan task.Task, numWorkers),
		CloseHandle: make(chan bool),
	}
}

func (p *ThreadPool) SubmitTaskCommon(task task.Task) error {
    if len(p.taskQueue) == p.queueSize {
        fmt.Println("Очередь заполнена, не удается добавить задачу")
        return ErrQueueFull
    }
    p.taskQueue <- task
    fmt.Println("Задача добавлена в очередь")
    return nil
}


func (p *ThreadPool) SubmitTaskPeriodoc(task task.Task, period time.Duration) {
	ticker := time.NewTicker(period)

	defer ticker.Stop()
	
	endTime := time.Now().Add(100 * time.Second) //  в теченеи 100 секунд будем добавлять 

	for {
		select {
		case <-ticker.C:
			if len(p.taskQueue) == p.queueSize {
				fmt.Println("Очередь заполнена, не удается добавить задачу")
				return
			}
			p.taskQueue <- task
			fmt.Println("Задача добавлена в очередь с периодом", period)
		case <- time.After(time.Until(endTime)):
			return
		}
	}
}
func (p *ThreadPool) SubmitTaskInterval(task task.Task, interval time.Duration) error {
    if len(p.taskQueue) == p.queueSize {
        fmt.Println("Очередь заполнена, не удается добавить задачу")
        return ErrQueueFull
    }

	time.Sleep(interval)

    p.taskQueue <- task
    fmt.Println("Задача добавлена в очередь")
    return nil
}

func (p *ThreadPool) Close() {
	close(p.CloseHandle)
	close(p.taskQueue)
	close(p.WorkerPool)
}

func (p *ThreadPool) Dispatch() {
    for {
        select {
        case task, ok := <-p.taskQueue:
            if !ok {
                return
            }

            jobChannel, ok := <-p.WorkerPool
            if !ok {
                fmt.Println("Нет свободных воркеров")
                return
            }

            select {
            case jobChannel <- task:
            case <-p.CloseHandle:
                return
            }
        case <-p.CloseHandle:
            return
        }
    }
}
