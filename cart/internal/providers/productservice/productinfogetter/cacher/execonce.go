package cacher

import (
	"sync"
	"sync/atomic"
)

type funcToBeExecutedOnce func() CacheValue

type executingOnce struct {
	key        CacheKey
	val        CacheValue
	once       sync.Once
	f          funcToBeExecutedOnce
	clientsNum atomic.Int64
	container  *execOnceCoordinator
}

func (eo *executingOnce) Execute() CacheValue {
	execWithSaving := func() {
		eo.val = eo.f()
	}
	eo.once.Do(execWithSaving)
	return eo.val
}

func (eo *executingOnce) Close() {
	eo.clientsNum.Add(-1)
	if eo.clientsNum.Load() == 0 {
		eo.container.deleteExecutor(eo)
	}
}

type execOnceCoordinator struct {
	mu   sync.Mutex
	data map[CacheKey]*executingOnce
}

func newExecOnceCoordinator() *execOnceCoordinator {
	return &execOnceCoordinator{data: make(map[CacheKey]*executingOnce)}
}

func (c *execOnceCoordinator) getExecutor(k CacheKey, f funcToBeExecutedOnce) *executingOnce {
	c.mu.Lock()
	defer c.mu.Unlock()
	exec, exist := c.data[k]
	if !exist {
		exec = &executingOnce{key: k, container: c, f: f}
		c.data[k] = exec
	}
	exec.clientsNum.Add(1)
	return exec
}

func (c *execOnceCoordinator) deleteExecutor(eo *executingOnce) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if eo.clientsNum.Load() == 0 {
		delete(c.data, eo.key)
	}
}
