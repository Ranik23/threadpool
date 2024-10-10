package worker

import (
	"threadpool/internal/task"
)

type Worker struct {
	jobChannel chan task.Task
	workerPool chan chan task.Task
	closeSIG   chan bool
}

func NewWorker(workerPool chan chan task.Task, closeSIG chan bool) *Worker {
	return &Worker{workerPool: workerPool, jobChannel: make(chan task.Task), closeSIG: closeSIG}
}

func (w *Worker) execute(task task.Task) {
	task.Run()
}


func (w *Worker) Start() {
	for {
		w.workerPool <- w.jobChannel
		select {
		case job := <-w.jobChannel:
			w.execute(job)
		case <-w.closeSIG:
			return
		}
	}
}