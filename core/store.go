package core

import (
	"sync"
	"time"
)

// Item wraps the raw byte data with an expiration timestamp.
type Item struct {
	Data      []byte
	ExpiresAt int64 // Unix nanosecond timestamp. 0 means it lives forever.
}

// Store represents the thread-safe, concurrent memory engine.
type Store struct {
	mu   sync.RWMutex
	data map[string]Item
}

// NewStore initializes a ready-to-use memory store.
func NewStore() *Store {
	return &Store{
		data: make(map[string]Item),
	}
}

// Set safely locks the engine and mutates the state.
// ttl is a time.Duration (e.g., 5 * time.Second). Pass 0 for no expiration.
func (s *Store) Set(key string, value []byte, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt int64
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl).UnixNano()
	}

	s.data[key] = Item{
		Data:      value,
		ExpiresAt: expiresAt,
	}
}

// Get retrieves the state. It handles "Passive Expiration" by checking
// the clock before returning the data to the user.
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	item, exists := s.data[key]
	s.mu.RUnlock() // CRITICAL: Unlock immediately before doing anything else to avoid deadlocks!

	if !exists {
		return nil, false
	}

	// If the item has a TTL and the current time has passed the expiration time
	if item.ExpiresAt > 0 && time.Now().UnixNano() > item.ExpiresAt {
		s.Delete(key) // We safely call Delete because we released the RLock above.
		return nil, false
	}

	return item.Data, true
}

// Delete safely removes a key from the map.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
// StartSweeper launches a background goroutine that periodically
// scans and removes expired keys to prevent memory leaks.
func (s *Store) StartSweeper(interval time.Duration) {
	// The 'go' keyword spins this infinite loop off into its own concurrent thread
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C // Blocks until the ticker fires
			s.sweep()
		}
	}()
}

// sweep locks the database and deletes dead keys.
func (s *Store) sweep() {
	s.mu.Lock() // We need a Write Lock because we are deleting from the map
	defer s.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range s.data {
		if item.ExpiresAt > 0 && now > item.ExpiresAt {
			delete(s.data, key)
		}
	}
}