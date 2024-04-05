package ratelimiterhttp

import "time"

type SystemTimeTicker struct {
	defaultTicker        *time.Ticker
	intervalBetweenTicks time.Duration
	ticksCh              chan struct{}
}

func NewSystemTimeTicker(rps int64) *SystemTimeTicker {
	interval := time.Second / time.Duration(rps)
	return &SystemTimeTicker{intervalBetweenTicks: interval}
}

func (s *SystemTimeTicker) Stop() {
	s.defaultTicker.Stop()
}

func (s *SystemTimeTicker) Start() {
	s.defaultTicker = time.NewTicker(s.intervalBetweenTicks)
	s.setupCh()
}

func (s *SystemTimeTicker) GetTickCh() <-chan struct{} {
	return s.ticksCh
}

func (s *SystemTimeTicker) setupCh() {
	s.ticksCh = make(chan struct{})
	go func() {
		for range s.defaultTicker.C {
			s.ticksCh <- struct{}{}
		}
	}()
}
