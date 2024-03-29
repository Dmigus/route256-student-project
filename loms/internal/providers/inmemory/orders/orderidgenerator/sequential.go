// Package orderidgenerator содержит генератор уникальных id заказов для использования InMemory хранилищем заказов
package orderidgenerator

import "sync"

// SequentialGenerator генерирует id последовательно
type SequentialGenerator struct {
	mu     sync.Mutex
	prevID int64
}

// NewSequentialGenerator создаёт новый SequentialGenerator, в котором при первом вызове NewID будет занчение startVal
func NewSequentialGenerator(startVal int64) *SequentialGenerator {
	return &SequentialGenerator{
		prevID: startVal - 1,
	}
}

// NewID генерирует и возращет новое значение id
func (s *SequentialGenerator) NewID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	newID := s.prevID + 1
	s.prevID = newID
	return newID
}
