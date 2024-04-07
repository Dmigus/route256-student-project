package ratelimiter

import "time"

// SystemTimeTicker это тикер, привязанный к системному времени. Генерирует тики равномерно по времени.
type SystemTimeTicker struct {
	defaultTicker        *time.Ticker
	intervalBetweenTicks time.Duration
	ticksCh              chan struct{}
}

// NewSystemTimeTicker создаёт SystemTimeTicker с таким интервалом, чтобы было rps тиков в секунду
func NewSystemTimeTicker(rps int64) *SystemTimeTicker {
	interval := time.Second / time.Duration(rps)
	return &SystemTimeTicker{intervalBetweenTicks: interval}
}

// Stop останавливает генерацию "тиков"
func (s *SystemTimeTicker) Stop() {
	s.defaultTicker.Stop()
}

// Start запускает новую генерацию "тиков". При этом старый канал больше не используется. Клиент должен вызвать GetTickCh для полуения актульаного канала.
func (s *SystemTimeTicker) Start() {
	s.defaultTicker = time.NewTicker(s.intervalBetweenTicks)
	s.setupCh()
}

// GetTickCh отдаёт актуальный канал, в который происходят тики (если происходят:))
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
