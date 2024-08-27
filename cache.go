package cache

import (
	"container/list"
	"sync"
)

type cacheItem struct {
	key   string
	value interface{}
}

type Cache struct {
	capacity int
	data     map[string]*list.Element
	mu       sync.RWMutex
	list     *list.List
}

func NewCache(size int) *Cache {
	return &Cache{
		capacity: size,
		data:     make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if the key already exists, update the value and move it to the front.
	if element, found := c.data[key]; found {
		c.list.MoveToFront(element)
		element.Value.(*cacheItem).value = value
		return
	}

	// check for capcity and do eviction if needed
	if c.list.Len() == c.capacity {
		backElement := c.list.Back()
		if backElement != nil {
			c.list.Remove(backElement)
			delete(c.data, backElement.Value.(*cacheItem).key)
		}
	}

	newElement := c.list.PushFront(&cacheItem{key, value})
	c.data[key] = newElement
}

func (c *Cache) Get(key string) any {
	if element, found := c.data[key]; found {
		c.list.MoveToFront(element)
		return element.Value.(*cacheItem).value
	}
	return nil
}
