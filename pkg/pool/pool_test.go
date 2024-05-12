package pool

import (
	"sync"
	"testing"
)

func TestPool_Go(t *testing.T) {
	opts := PoolOptions{
		Max: 3,
	}
	p := New(&opts)
	var wg sync.WaitGroup
	wg.Add(opts.Max)

	for i := 0; i < opts.Max; i++ {
		p.Go(func() {
			wg.Done()
		})
	}

	wg.Wait()
}

func TestPool_GoWithMoreTasksThanMax(t *testing.T) {
	opts := PoolOptions{
		Max: 3,
	}
	p := New(&opts)
	var wg sync.WaitGroup
	taskCount := 10
	wg.Add(taskCount)

	for i := 0; i < taskCount; i++ {
		p.Go(func() {
			wg.Done()
		})
	}

	wg.Wait()
}
