package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int64) error {
	taskChan := make(chan Task)
	var errCnt int64 = 0
	errChan := make(chan struct{})
	stopChan := make(chan struct{})
	var wg sync.WaitGroup

	// ===========================================
	// с runner вроде все ок
	runner := func() {
		defer wg.Done()
		for {
			select {
			case curren_task := <-taskChan:
				err := current_task()
				if err {
					// atomic.AddInt32(errCnt)
					errChan <- struct{}{}
				}
			case <-stopChan:
				return
			}

		}
	}
	// ===========================================

	for i := 0; i < n; i++ {
		wg.Add(1)
		go runner()
	}

	for {
		select {
		case <-errChan:
			atomic.AddInt64(&errCnt, 1)
			if errCnt >= m {
				close(stopChan)
				break
			}
		}
	}

	return nil
}
