package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	pool := WorkersPool{
		poolSize:    n,
		maxErrCount: m,
		tasks:       tasks,
		mu:          &sync.Mutex{},
	}
	err := pool.RunTasks()
	return err
}

// WorkersPool - структура для работы с пул воркерами.
type WorkersPool struct {
	poolSize       int
	maxErrCount    int
	tasks          []Task
	mu             *sync.Mutex
	completedTasks int
	errCount       int
}

// RunTasks запускает выполнение задач.
func (w *WorkersPool) RunTasks() error {
	var wg sync.WaitGroup
	taskChan := make(chan Task, len(w.tasks))
	statusChan := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())

	go w.taskMonitor(cancel, statusChan)

	wg.Add(w.poolSize)
	for i := 0; i < w.poolSize; i++ {
		go func() {
			defer wg.Done()
			w.startWorker(taskChan, statusChan)
		}()
	}

	for _, task := range w.tasks {
		taskChan <- task
	}

	<-ctx.Done()
	close(taskChan)
	wg.Wait()
	close(statusChan)
	if !w.errLimitCheck() {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (w *WorkersPool) startWorker(in <-chan Task, status chan<- bool) {
	for task := range in {
		if w.errLimitCheck() {
			err := task()
			status <- err == nil
		}
	}
}

func (w *WorkersPool) taskMonitor(cancel context.CancelFunc, status <-chan bool) {
	for ok := range status {
		w.completedTasks++
		if !ok {
			w.mu.Lock()
			w.errCount++
			w.mu.Unlock()
			if !w.errLimitCheck() {
				cancel()
			}
		}
		if !w.taskDoneCheck() {
			cancel()
		}
	}
}

func (w *WorkersPool) errLimitCheck() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.maxErrCount <= 0 { // если лимит ошибок <= 0, то игнорируем ошибки в принципе
		return true
	}
	if w.errCount >= w.maxErrCount {
		return false
	}
	return true
}

func (w *WorkersPool) taskDoneCheck() bool {
	return w.completedTasks != len(w.tasks)
}
