package cache

import (
	"strconv"
	"sync"
	"testing"
)

func TestBasicOperations(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("B", 2)

	if val := cache.Get("A"); val != 1 {
		t.Errorf("Expected Get(A) to return 1, got %d", val.(int))
	}

	if val := cache.Get("B"); val != 2 {
		t.Errorf("Expected Get(B) to return 2, got %d", val.(int))
	}
}

type custom struct {
	Item string
}

func TestCustomStruct(t *testing.T) {
	cache := NewCache(2)

	customVals := []custom{
		{Item: "item1"},
		{Item: "item2"},
	}

	cache.Set("A", customVals[0])
	cache.Set("B", customVals[1])

	if val := cache.Get("A"); val.(custom) != customVals[0] {
		t.Errorf("Expected Get(A) to return item A %v, got %v", customVals[0], val)
	}

	if val := cache.Get("B"); val.(custom) != customVals[1] {
		t.Errorf("Expected Get(B) to return 2 %v, got %v", customVals[0], val)
	}
}

func TestUpdateValue(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("A", 10)

	if val := cache.Get("A"); val != 10 {
		t.Errorf("Expected Get(A) to return 10 after update, got %d", val.(int))
	}
}

func TestCapacityOverflow(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("B", 2)
	cache.Set("C", 3) // evicts key A

	if val := cache.Get("A"); val != nil {
		t.Errorf("Expected Get(A) to return nil (evicted), got %d", val.(int))
	}

	if val := cache.Get("B"); val != 2 {
		t.Errorf("Expected Get(B) to return 2, got %d", val.(int))
	}

	if val := cache.Get("C"); val != 3 {
		t.Errorf("Expected Get(C) to return 3, got %d", val.(int))
	}
}

func TestGetUpdatesRecency(t *testing.T) {
	cache := NewCache(2)

	cache.Set("A", 1)
	cache.Set("B", 2)
	cache.Get("A")    // updates recency of key A
	cache.Set("C", 3) // evicts key B

	if val := cache.Get("B"); val != nil {
		t.Errorf("Expected Get(B) to return nil (evicted), got %d", val.(int))
	}

	if val := cache.Get("A"); val != 1 {
		t.Errorf("Expected Get(A) to return 1, got %d", val.(int))
	}

	if val := cache.Get("C"); val != 3 {
		t.Errorf("Expected Get(C) to return 3, got %d", val.(int))
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewCache(100)
	var wg sync.WaitGroup

	// Writing concurrently
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(k, v int) {
			defer wg.Done()
			cache.Set(strconv.Itoa(i), v)
		}(i, i*i)
	}

	// Reading concurrently
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			_ = cache.Get(strconv.Itoa(i))
		}(i)
	}

	wg.Wait()

	// Verify some values
	for i := 900; i < 1000; i++ {
		val := cache.Get(strconv.Itoa(i))
		if val != i*i && val != nil {
			t.Errorf("Expected Get(%d) to return %d or nil, got %d", i, i*i, val.(int))
		}
	}
}

func TestNegativeCapacity(t *testing.T) {
	cache := NewCache(-1)

	cache.Set("A", 1)

	if val := cache.Get("A"); val != nil {
		t.Errorf("Expected Get(A) to return nil for negative capacity cache, got %d", val.(int))
	}
}

func TestZeroCapacity(t *testing.T) {
	cache := NewCache(0)

	cache.Set("A", 1)

	if val := cache.Get("A"); val != nil {
		t.Errorf("Expected Get(A) to return nil for negative capacity cache, got %d", val.(int))
	}
}
