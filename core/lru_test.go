package core

import (
	"testing"
)

func TestLRU_BasicPutAndGet(t *testing.T) {
	cache := NewLRU(2) // Capacity of 2 for easy testing

	cache.Put("user_1", "Adityajyoti")
	
	val, found := cache.Get("user_1")
	if !found || val != "Adityajyoti" {
		t.Fatalf("Expected 'Adityajyoti', got '%v'", val)
	}
}

func TestLRU_Eviction(t *testing.T) {
	// Initialize cache with a strict capacity of 3
	cache := NewLRU(3)

	// Fill the cache
	cache.Put("key1", "val1")
	cache.Put("key2", "val2")
	cache.Put("key3", "val3")

	// Access key1 so it becomes the MOST recently used
	cache.Get("key1")

	// Insert a 4th key. This MUST trigger an eviction.
	// Since key1 was just accessed, key2 is now the LEAST recently used and should die.
	cache.Put("key4", "val4")

	// Verify key2 was evicted
	_, found := cache.Get("key2")
	if found {
		t.Fatal("LRU eviction failed: key2 should have been evicted to prevent OOM")
	}

	// Verify key1 survived because we accessed it
	_, found = cache.Get("key1")
	if !found {
		t.Fatal("LRU eviction failed: key1 was accessed and should have survived")
	}

	// Verify the new key4 was safely added
	_, found = cache.Get("key4")
	if !found {
		t.Fatal("LRU eviction failed: key4 was not found after insertion")
	}
}