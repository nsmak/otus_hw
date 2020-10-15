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
	pool := PoolOfWorkers{
		poolSize:    n,
		maxErrCount: m,
		tasks:       tasks,
		mu:          &sync.Mutex{},
	}
	err := pool.RunTasks()
	return err
}

// PoolOfWorkers - структура для работы с пул воркерами.
type PoolOfWorkers struct {
	poolSize       int
	maxErrCount    int
	tasks          []Task
	mu             *sync.Mutex
	completedTasks int
	errCount       int
}

// RunTasks запускает выполнение задач.
func (p *PoolOfWorkers) RunTasks() error {
	var wg sync.WaitGroup
	taskChan := make(chan Task, len(p.tasks))
	statusChan := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())

	go p.taskMonitor(cancel, statusChan)

	for i := 0; i < p.poolSize; i++ {
		wg.Add(1)
		go func() {
			p.startWorker(taskChan, statusChan)
			wg.Done()
		}()
	}

	for _, task := range p.tasks {
		taskChan <- task
	}

	<-ctx.Done()
	close(taskChan)
	wg.Wait()
	close(statusChan)
	if !p.errLimitCheck() {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (p *PoolOfWorkers) startWorker(in <-chan Task, status chan<- bool) {
	for task := range in {
		if p.errLimitCheck() {
			err := task()
			status <- err == nil
		}
	}
}

func (p *PoolOfWorkers) taskMonitor(cancel context.CancelFunc, status <-chan bool) {
	for ok := range status {
		p.mu.Lock()
		p.completedTasks++
		p.mu.Unlock()
		if !ok {
			p.mu.Lock()
			p.errCount++
			p.mu.Unlock()
			if !p.errLimitCheck() {
				cancel()
			}
		}
		if !p.taskDoneCheck() {
			cancel()
		}
	}
}

func (p *PoolOfWorkers) errLimitCheck() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.maxErrCount <= 0 { // если лимит ошибок <= 0, то игнорируем ошибки в принципе
		return true
	}
	if p.errCount >= p.maxErrCount {
		return false
	}
	return true
}

func (p *PoolOfWorkers) taskDoneCheck() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.completedTasks != len(p.tasks)
}
