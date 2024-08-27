package cache

import "testing"

func TestBasicCache(t *testing.T) {
	cache := New(2).LRU().Build()

	cache.Set("A", 1)
	cache.Set("B", 2)

	if val := cache.Get("A"); val != 1 {
		t.Errorf("Expected Get(A) to return 1, got %d", val.(int))
	}

	if val := cache.Get("B"); val != 2 {
		t.Errorf("Expected Get(B) to return 2, got %d", val.(int))
	}
}
