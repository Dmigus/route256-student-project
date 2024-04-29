// Package inmemorycache содержит in-memory реализацию кэша
package inmemorycache

import (
	"container/list"
	"context"
	"sync"

	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter/cacher"
)

const defaultMaxSize = uint(10)

var zeroVal cacher.CacheValue

type (
	queueElemType struct {
		k cacher.CacheKey
		v cacher.CacheValue
	}
	// InMemoryCache это in-memory реализация кэша
	InMemoryCache struct {
		mu      sync.Mutex
		kv      map[cacher.CacheKey]*list.Element
		queue   *list.List
		maxSize uint
	}
)

// NewInMemoryCache создаёт новый InMemoryCache
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

// Get возвращает значение, сохранённое в кэше и true, если оно там было. Если не было, то false
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

// Store сохраняет новое или обновляет старое значение в кэше. Вытесняет наименее актуальное значение, если это необходимо.
func (c *InMemoryCache) Store(_ context.Context, k cacher.CacheKey, v cacher.CacheValue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	valElement, present := c.kv[k]
	if present {
		c.queue.Remove(valElement)
		delete(c.kv, k)
	} else if c.Size() >= c.maxSize {
		c.deleteLastLocked()
	}
	c.insertLocked(k, v)
}

// Size возвращает текущее количество элементов в кэше
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
