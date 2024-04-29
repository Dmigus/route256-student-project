package cacher

import (
	"sync"
	"sync/atomic"
)

type funcToBeExecutedAtMostOnce func() CacheValue

// oneTimeExecutor это одноразовый executor, который должен выполнить переданную ему функцию, сохранить себе результат
type oneTimeExecutor struct {
	key        CacheKey
	val        CacheValue
	once       sync.Once
	clientsNum atomic.Int64
	container  *execOnceCoordinator
}

func (eo *oneTimeExecutor) Execute(f funcToBeExecutedAtMostOnce) CacheValue {
	defer eo.close()
	eo.once.Do(func() {
		eo.val = f()
	})
	return eo.val
}

func (eo *oneTimeExecutor) close() {
	eo.clientsNum.Add(-1)
	if eo.clientsNum.Load() == 0 {
		eo.container.deleteExecutor(eo)
	}
}

type execOnceCoordinator struct {
	mu   sync.Mutex
	data map[CacheKey]*oneTimeExecutor
}

func newExecOnceCoordinator() *execOnceCoordinator {
	return &execOnceCoordinator{data: make(map[CacheKey]*oneTimeExecutor)}
}

func (c *execOnceCoordinator) getExecutor(k CacheKey) *oneTimeExecutor {
	c.mu.Lock()
	defer c.mu.Unlock()
	exec, exist := c.data[k]
	if !exist {
		exec = &oneTimeExecutor{key: k, container: c}
		c.data[k] = exec
	}
	exec.clientsNum.Add(1)
	return exec
}

func (c *execOnceCoordinator) deleteExecutor(eo *oneTimeExecutor) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if eo.clientsNum.Load() == 0 {
		delete(c.data, eo.key)
	}
}
