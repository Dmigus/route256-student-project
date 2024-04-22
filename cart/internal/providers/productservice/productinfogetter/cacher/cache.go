package cacher

import (
	"container/list"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"
	"sync"
)

const defaultMaxSize = uint(10)

type (
	key struct {
		method  string
		request productinfogetter.GetProductRequest
	}
	value struct {
		response productinfogetter.GetProductResponse
		err      error
	}
	queueElemType struct {
		k key
		v value
	}
)

var zeroVal value

type cache struct {
	mu      sync.Mutex
	kv      map[key]*list.Element
	queue   *list.List
	maxSize uint
}

func newCache() *cache {
	return &cache{
		kv:      make(map[key]*list.Element),
		queue:   list.New(),
		maxSize: defaultMaxSize,
	}
}

func (c *cache) Get(k key) (value, bool) {
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

func (c *cache) Insert(k key, v value) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Size() >= c.maxSize {
		c.deleteLastLocked()
	}
	c.insertLocked(k, v)
}

func (c *cache) Size() uint {
	return uint(len(c.kv))
}

func (c *cache) deleteLastLocked() {
	outdated := c.queue.Back()
	outdatedKey := outdated.Value.(queueElemType).k
	c.queue.Remove(outdated)
	delete(c.kv, outdatedKey)
}

func (c *cache) insertLocked(k key, v value) {
	queueElemVal := queueElemType{k: k, v: v}
	elem := c.queue.PushFront(queueElemVal)
	c.kv[k] = elem
}
