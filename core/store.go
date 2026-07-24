package core

import (
	"sync"
	"time"
)

// Store is the thread-safe engine wrapping the LRU cache
type Store struct {
	mu  sync.Mutex
	lru *LRUCache
}

// NewStore initializes a concurrent-safe key-value store
func NewStore(capacity int) *Store {
	return &Store{
		lru: NewLRU(capacity),
	}
}

// Put safely acquires a lock and inserts/updates a key (Used by embedded API)
func (s *Store) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lru.Put(key, value)
}

// Set handles TTL and byte-slice mutations for the WAL and TCP server
func (s *Store) Set(key string, value []byte, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Storing as string in the LRU. TTL tracking can be expanded here later.
	s.lru.Put(key, string(value)) 
}

// Get safely acquires a lock and retrieves a value
func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lru.Get(key)
}

// Delete removes a key from the store safely in O(1) time
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if node, found := s.lru.cache[key]; found {
		s.lru.removeNode(node)
		delete(s.lru.cache, key)
	}
}

// StartSweeper handles background expiration (stubbed for main.go compliance)
func (s *Store) StartSweeper(interval time.Duration) {
	// Background ticker for TTL cleanup to be implemented
}
