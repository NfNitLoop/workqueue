package workqueue

import (
	"sync"
)

// Run runs `concurrency` goroutines to work items submitted to `workQueue`.
// 
// When Run(...) is finished, all work is guaranteed to have completed.
func Run(concurrency int, callback func(workQueue WorkQueue)) {
	workQueue := &workQueue {
		channel: make(chan func()),
	}

	if concurrency < 1 {
		concurrency = 1
	}

	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		// Start worker goroutines:
		go func() {
			defer wg.Done()
			for callback := range workQueue.channel {
				callback()
			}
		}()
	}

	// Let the user send workers to the queue:
	callback(workQueue)
	// Tell workers no more items are coming.
	close(workQueue.channel)
	// Make sure all of our workers finish:
	wg.Wait()
}

// WorkQueue runs your jobs.
type WorkQueue interface {
	// Submit submits a job to the WorkQueue to run.
	Submit(job func())
}

// implements WorkQueue
type workQueue struct {
	channel chan func()
}

func (w *workQueue) Submit(job func()) {
	w.channel <- job
}
