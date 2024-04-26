package inmemorycache

import (
	"container/list"
	"context"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter/cacher"
	"sync"
)

const defaultMaxSize = uint(10)

var zeroVal cacher.CacheValue

type (
	queueElemType struct {
		k cacher.CacheKey
		v cacher.CacheValue
	}
	InMemoryCache struct {
		mu      sync.Mutex
		kv      map[cacher.CacheKey]*list.Element
		queue   *list.List
		maxSize uint
	}
)

func NewInMemoryCache(opts ...Option) *InMemoryCache {
	cache := &InMemoryCache{
		kv:      make(map[cacher.CacheKey]*list.Element),
		queue:   list.New(),
		maxSize: defaultMaxSize,
	}
	for _, opt := range opts {
		opt.apply(cache)
	}
	return cache
}

func (c *InMemoryCache) Get(_ context.Context, k cacher.CacheKey) (cacher.CacheValue, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	valElement, present := c.kv[k]
	if !present {
		return zeroVal, false
	}
	c.queue.MoveToFront(valElement)
	val := valElement.Value.(queueElemType)
	return val.v, true
}

func (c *InMemoryCache) Store(_ context.Context, k cacher.CacheKey, v cacher.CacheValue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Size() >= c.maxSize {
		c.deleteLastLocked()
	}
	c.insertLocked(k, v)
}

func (c *InMemoryCache) Size() uint {
	return uint(len(c.kv))
}

func (c *InMemoryCache) deleteLastLocked() {
	outdated := c.queue.Back()
	outdatedKey := outdated.Value.(queueElemType).k
	c.queue.Remove(outdated)
	delete(c.kv, outdatedKey)
}

func (c *InMemoryCache) insertLocked(k cacher.CacheKey, v cacher.CacheValue) {
	queueElemVal := queueElemType{k: k, v: v}
	elem := c.queue.PushFront(queueElemVal)
	c.kv[k] = elem
}
