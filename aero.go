package aerokv

import (
	"github.com/GiTaDi-CrEaTe/aero-kv/core"
)

// DB represents a fully embedded, thread-safe aero-kv instance.
type DB struct {
	store *core.Store
}

// Open initializes a new embedded database with the given OOM-safe capacity.
func Open(capacity int) *DB {
	return &DB{
		store: core.NewStore(capacity),
	}
}

// Set stores a key-value pair safely across concurrent goroutines.
func (db *DB) Set(key, value string) {
	db.store.Put(key, value)
}

// Get retrieves a value by key, updating the internal LRU eviction state.
func (db *DB) Get(key string) (string, bool) {
	return db.store.Get(key)
}