package cache

import (
	"testing"
)

func TestBasicOperations(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("B", 2)

	if val := cache.Get("A"); val != 1 {
		t.Errorf("Expected Get(A) to return 1, got %d", val)
	}

	if val := cache.Get("B"); val != 2 {
		t.Errorf("Expected Get(B) to return 2, got %d", val)
	}
}

func TestUpdateValue(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("A", 10)

	if val := cache.Get("A"); val != 10 {
		t.Errorf("Expected Get(A) to return 10 after update, got %d", val)
	}
}

func TestCapacityOverflow(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("B", 2)
	cache.Set("C", 3) // evicts key A

	if val := cache.Get("A"); val != nil {
		t.Errorf("Expected Get(A) to return nil (evicted), got %d", val)
	}

	if val := cache.Get("B"); val != 2 {
		t.Errorf("Expected Get(B) to return 2, got %d", val)
	}

	if val := cache.Get("C"); val != 3 {
		t.Errorf("Expected Get(C) to return 3, got %d", val)
	}
}

func TestGetUpdatesRecency(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("B", 2)
	cache.Get("A")    // updates recency of key A
	cache.Set("C", 3) // evicts key B

	if val := cache.Get("B"); val != nil {
		t.Errorf("Expected Get(B) to return nil (evicted), got %d", val)
	}

	if val := cache.Get("A"); val != 1 {
		t.Errorf("Expected Get(A) to return 1, got %d", val)
	}

	if val := cache.Get("C"); val != 3 {
		t.Errorf("Expected Get(C) to return 3, got %d", val)
	}
}
