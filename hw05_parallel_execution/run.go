package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	if n <= 0 {
		return nil
	}

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var (
		errCnt int64
		wg     sync.WaitGroup
	)
	tasksCh := make(chan Task, len(tasks))

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				t, ok := <-tasksCh
				if !ok {
					return
				}

				if err := t(); err != nil {
					atomic.AddInt64(&errCnt, 1)
				}

				if atomic.LoadInt64(&errCnt) >= int64(m) {
					return
				}
			}
		}()
	}

	for i := 0; i < len(tasks); i++ {
		tasksCh <- tasks[i]
	}
	close(tasksCh)

	wg.Wait()
	if errCnt >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
