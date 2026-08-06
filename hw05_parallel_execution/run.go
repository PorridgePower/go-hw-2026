package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChan := make(chan Task)
	errChan := make(chan struct{})
	stopChan := make(chan struct{})
	doneChan := make(chan struct{})
	var errCnt int32 = 0
	var wg sync.WaitGroup

	// ===========================================
	// с runner вроде все ок
	runner := func() {
		defer wg.Done()
		for {
			select {
			case currentTask, ok := <-taskChan:
				if !ok {
					return
				}
				err := currentTask()
				if err != nil {
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

	wg.Add(1)
	// Coordinator
	go func() {
		defer wg.Done()
		// defer close(taskChan)
		// defer close(taskChan)
		for _, t := range tasks {
			select {
			case <-stopChan:
				// close(taskChan)
				return
			case taskChan <- t:
				continue
			}
		}
		doneChan <- struct{}{}

	}()

	for {
		select {
		case <-errChan:
			atomic.AddInt32(&errCnt, 1)
			if m > 0 && errCnt >= int32(m) {
				close(stopChan)
				close(taskChan)
				return ErrErrorsLimitExceeded
			}
		// here check af all task processed
		case <-doneChan:
			// break
			wg.Wait()
			return nil
		}
	}
	wg.Wait()
	return nil
}
