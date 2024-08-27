package cache

const (
	TYPE_LRU = "lru"
)

type Cache interface {
	Set(key string, value any)
	Get(key string) any
	Purge()
}

type CacheBuilder struct {
	capacity int
	ctype    string
}

func New(capacity int) *CacheBuilder {
	return &CacheBuilder{capacity: capacity}
}

func (cb *CacheBuilder) Build() Cache {
	if cb.capacity <= 0 {
		panic("cache size <= 0")
	}

	switch cb.ctype {
	case TYPE_LRU:
		return newLRUCache(cb.capacity)
	default:
		panic("unknown type " + cb.ctype)
	}

}

func (cb *CacheBuilder) LRU() *CacheBuilder {
	cb.ctype = TYPE_LRU
	return cb
}
