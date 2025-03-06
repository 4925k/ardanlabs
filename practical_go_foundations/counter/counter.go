package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// var mu sync.Mutex
	// count := 0
	// ALTERNATIVE WITH ATOMIC
	var count int64

	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10_000 {
				// mu.Lock()
				// count++
				// mu.Unlock()

				atomic.AddInt64(&count, 1)
			}
		}()
	}

	wg.Wait()          // Wait for all goroutines to finish
	fmt.Println(count) // Should print 100000
}
