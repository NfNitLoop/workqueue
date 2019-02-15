package workqueue

import (
	"fmt"
	"image"
	"sync"
)


// If you're not using workqueue, here's the raw Go boilerplate you'd have to write:
func Example_a() {
	items := []image.Point {
		{1,1}, {2,2}, {3,3}, {4,4}, {5,5},
	}

	workQueue := make(chan *image.Point)
	concurrency := 2
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		// Start worker goroutines:
		go func() {
			defer wg.Done()
			for point := range workQueue {
				point.X = point.X * 2
				point.Y = point.Y * 2
			}
		}()
	}

	// Send workers to the queue:
	for i := range items {
		workQueue <- &items[i]
	}
	// Tell workers no more items are coming.
	close(workQueue)

	// Make sure all of our workers finish:
	wg.Wait()

	for _, item := range items {
		fmt.Println(item.String())
	}

	// Output:
	// (2,2)
	// (4,4)
	// (6,6)
	// (8,8)
	// (10,10)
}

// But there's a lot less boilerplate with workqueue:
func Example_b() {
	items := []image.Point {
		{1,1}, {2,2}, {3,3}, {4,4}, {5,5},
	}

	// A job is just a func(). Create a job to scale points:
	scale2 := func(p *image.Point) func() {
		return func() {
			p.X = p.X * 2
			p.Y = p.Y * 2
		}
	}

	Run(2, func(queue WorkQueue) {
		for i := range items {
			queue.Submit(scale2(&items[i]))
		}
	})

	for _, item := range items {
		fmt.Println(item.String())
	}

	// Output:
	// (2,2)
	// (4,4)
	// (6,6)
	// (8,8)
	// (10,10)
}

// You can still use channels for output if you like:
func Example_c() {
	items := []image.Point {
		{1,1}, {2,2}, {3,3}, {4,4}, {5,5},
	}

	output := make(chan image.Point)

	// This job sends output to our channel, instead of
	// modifying data in-place as in the previous example:
	scale2 := func(p image.Point) func() {
		return func() {
			output <- p.Mul(2)
		}
	}

	// Submitting work to the queue may block, so we do it in
	// a goroutine:
	go func() { 
		Run(2, func(queue WorkQueue) {
			for _, item := range items {
				queue.Submit(scale2(item))
			}	
		})

		// ... but once Run() finishes, all work has completed
		// so we can close the channel. No WaitGroup needed:
		close(output)
	}()

	// Print output as it happens:
	for point := range output {
		fmt.Println(point.String())
	}

	// Unordered output:
	// (2,2)
	// (4,4)
	// (6,6)
	// (8,8)
	// (10,10)
}

// Since Go will automatically convert yourStruct.method to a closure of type
// `func()`, using workqueue with methods is even easier:
func Example_d() {
	items := []*Foo {
		{1}, {2}, {3},
	}

	Run(2, func(queue WorkQueue) {
		for _, item := range items {
			queue.Submit(item.double)
		}
	})

	for _, item := range items {
		fmt.Println(item.value)
	}

	// Output:
	// 2
	// 4
	// 6
}

type Foo struct {
	value int
}
func (f *Foo) double() {
	f.value = f.value * 2
}