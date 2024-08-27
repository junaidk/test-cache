package cache

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	value  interface{}
	expiry time.Time // TTL for a key
}

type Cache struct {
	capacity int
	data     map[string]CacheItem
	mu       sync.RWMutex
	list     *list.List
}

func NewCache(size int) *Cache {
	return &Cache{
		capacity: size,
		data:     make(map[string]CacheItem),
		list:     list.New(),
	}
}
