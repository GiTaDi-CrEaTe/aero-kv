package core

import (
	"strconv"
	"sync"
	"testing"
)

// This benchmarks how fast your store handles massive concurrent writes
func BenchmarkStore_ConcurrentWrites(b *testing.B) {
	// Assuming your Store has a NewStore constructor and wraps the LRU
	store := NewStore(1000) 
	var wg sync.WaitGroup

	b.ResetTimer() // Start the clock only after initialization

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		// Launching concurrent goroutines for every single operation
		go func(n int) {
			defer wg.Done()
			store.Put("key_"+strconv.Itoa(n), "value")
		}(i)
	}
	
	wg.Wait()
}