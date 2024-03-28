package orderidgenerator

import "sync"

// SequentialGenerator генерирует id последовательно
type SequentialGenerator struct {
	mu     sync.Mutex
	prevID int64
}

// NewSequentialGenerator создаёт новый SequentialGenerator, в котором при первом вызове NewId будет занчение startVal
func NewSequentialGenerator(startVal int64) *SequentialGenerator {
	return &SequentialGenerator{
		prevID: startVal - 1,
	}
}

// NewId генерирует и возращет новое значение id
func (s *SequentialGenerator) NewId() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	newId := s.prevID + 1
	s.prevID = newId
	return newId
}
