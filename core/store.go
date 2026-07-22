package core

import "sync"

// Store is the thread-safe engine wrapping the LRU cache
type Store struct {
	mu  sync.Mutex // Strict Mutex to prevent linked-list corruption on reads
	lru *LRUCache
}

// NewStore initializes a concurrent-safe key-value store
func NewStore(capacity int) *Store {
	return &Store{
		lru: NewLRU(capacity),
	}
}

// Put safely acquires a lock and inserts/updates a key
func (s *Store) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lru.Put(key, value)
}

// Get safely acquires a full lock because an LRU read mutates the underlying list
func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lru.Get(key)
}