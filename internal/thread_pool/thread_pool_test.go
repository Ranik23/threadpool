package threadpool_test

import (
	"testing"
	"threadpool/internal/task"
	threadpool "threadpool/internal/thread_pool"
	"threadpool/internal/worker"
	"time"
)




func BenchmarkThreadPool_LowLoad(b *testing.B) {
    numWorkers := 2
    queueTask := 10

	signalEmptyQueue := make(chan bool)

    pool := threadpool.NewThreadPool(numWorkers, queueTask)

    for i := 0; i < numWorkers; i++ {
        go func() {
            w := worker.NewWorker(pool.WorkerPool, signalEmptyQueue)
            w.Start()
        }()
    }

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        for j := 0; j < queueTask; j++ {
            index := j
            t := task.TTask{Id: index}
            pool.SubmitTaskCommon(t)
        }
        
        go pool.Dispatch()
    }

	b.StopTimer()

	time.Sleep(5 * time.Second)

	pool.Close()
}

func BenchmarkThreadPool_HighLoad(b *testing.B) {
	numWorkers := 10
	queueTask := 100000

	signalEmptyQueue := make(chan bool)

	pool := threadpool.NewThreadPool(numWorkers, queueTask)

	for i := 0; i < numWorkers; i++ {
		go func() {
			w := worker.NewWorker(pool.WorkerPool, signalEmptyQueue)
			w.Start()
		}()
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < queueTask; j++ {
			index := j
			t := task.TTask{Id: index}
			pool.SubmitTaskCommon(t)
		}
		go pool.Dispatch()
	}

	b.StopTimer()

	time.Sleep(100 * time.Second)

	pool.Close()
}
