package orderidgenerator

import "sync"

type SequentialGenerator struct {
	mu     sync.Mutex
	prevID int64
}

func NewSequentialGenerator() *SequentialGenerator {
	return &SequentialGenerator{
		prevID: 0,
	}
}

func (s *SequentialGenerator) NewId() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	newId := s.prevID + 1
	s.prevID = newId
	return newId
}
