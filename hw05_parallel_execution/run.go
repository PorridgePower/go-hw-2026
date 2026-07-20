package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func checkLimit(x, limit int64) error {
	switch {
	case limit > 0 && x >= limit:
		return ErrErrorsLimitExceeded
	case limit == 0 && x > limit:
		return ErrErrorsLimitExceeded
	}
	return nil
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return nil
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	tasksCh := make(chan Task)
	stopCh := make(chan struct{})
	var errCount atomic.Int64
	var wg sync.WaitGroup
	var once sync.Once

	stop := func() {
		once.Do(func() { close(stopCh) })
	}

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			case task, ok := <-tasksCh:
				if !ok {
					return
				}
				if err := task(); err != nil {
					errCount.Add(1)
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		defer close(tasksCh)
		for _, t := range tasks {
			select {
			case <-stopCh:
				return
			case tasksCh <- t:
			}
		}
	}()

	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				if err := checkLimit(errCount.Load(), int64(m)); err != nil {
					stop()
					return
				}
			}
		}
	}()

	wg.Wait()
	stop()

	return checkLimit(errCount.Load(), int64(m))
}
