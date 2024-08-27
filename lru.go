package cache

import (
	"container/list"
	"sync"
)

type lruItem struct {
	key   string
	value interface{}
}

type lruCache struct {
	capacity int
	data     map[string]*list.Element
	mu       sync.RWMutex
	list     *list.List
}

func newLRUCache(size int) Cache {
	return &lruCache{
		capacity: size,
		data:     make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *lruCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if the key already exists, update the value and move it to the front.
	if element, found := c.data[key]; found {
		c.list.MoveToFront(element)
		element.Value.(*lruItem).value = value
		return
	}

	// check for capcity and do eviction if needed
	if c.list.Len() == c.capacity {
		backElement := c.list.Back()
		if backElement != nil {
			c.list.Remove(backElement)
			delete(c.data, backElement.Value.(*lruItem).key)
		}
	}

	newElement := c.list.PushFront(&lruItem{key, value})
	c.data[key] = newElement
}

func (c *lruCache) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if element, found := c.data[key]; found {
		c.list.MoveToFront(element)
		return element.Value.(*lruItem).value
	}
	return nil
}

func (c *lruCache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*list.Element)
	c.list.Init()
}
