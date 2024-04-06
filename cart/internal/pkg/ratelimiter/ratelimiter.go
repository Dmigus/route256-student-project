package ratelimiter

import (
	"context"
	"sync"
	"sync/atomic"
)

type ticker interface {
	Stop()
	Start()
	GetTickCh() <-chan struct{}
}

// RateLimiter это структура, которая позволяет ограничивать количество используемых ресурсов в единицу времени
type RateLimiter struct {
	waiting    chan struct{}
	capacity   uint64
	waitingCnt atomic.Int64

	mu            sync.Mutex
	availableInt  uint64
	ticker        ticker
	tickerRunning bool
}

// NewRateLimiter возвращает новый *RateLimiter, готовый к запуску.
// Сразу после инициализации количество ресурсов полное, то есть первые capacity запросов Acquire будут удовлетворены без блокировки.
// После запуска ресурс будет пополняться переданным тикером.
func NewRateLimiter(capacity uint64, ticker ticker) *RateLimiter {
	rl := &RateLimiter{availableInt: capacity, capacity: capacity, ticker: ticker, waiting: make(chan struct{})}
	return rl
}

func (lb *RateLimiter) addItemLocked() {
	// сначала освобождаем ожидающих
	select {
	case <-lb.waiting:
		return
	default:
	}
	// если ожидающих не было
	if lb.availableInt < lb.capacity {
		lb.availableInt++
	}
	return
}

// Acquire это блокирующий вызов, который завершается в случаях:
// 1) если получилось использовать ресурс. В таком случае возращаемое значение nil
// 2) если завершился переданный контекст. В таком  случае возвращается прчина завершения контекста
func (lb *RateLimiter) Acquire(ctx context.Context) error {
	lb.mu.Lock()
	if lb.availableInt > 0 {
		lb.availableInt--
		lb.startTickerIfStoppedLocked()
		lb.mu.Unlock()
		return nil
	}
	// увеличиваем счётчик ожидающих до запуска тикера
	lb.waitingCnt.Add(1)
	// после окончания блокировки уменьшаем счётчик
	defer lb.waitingCnt.Add(-1)
	lb.startTickerIfStoppedLocked()
	lb.mu.Unlock()
	return lb.waitForAvailableTick(ctx)
}

func (lb *RateLimiter) waitForAvailableTick(ctx context.Context) error {
	select {
	case lb.waiting <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// запуск заполнения доступного ресурса. Неблокирующий.
func (lb *RateLimiter) startTickerIfStoppedLocked() {
	if lb.tickerRunning {
		return
	}
	lb.ticker.Start()
	lb.tickerRunning = true
	go lb.ticking()
}

func (lb *RateLimiter) ticking() {
	for {
		select {
		case <-lb.ticker.GetTickCh():
			lb.mu.Lock()
			lb.addItemLocked()
			// если полностью заполнен capacity и нет ожидающих, то останавливаем тикер
			if lb.couldTickerBeStoppedLocked() {
				lb.tickerRunning = false
				lb.ticker.Stop()
				lb.mu.Unlock()
				return
			}
			lb.mu.Unlock()
		}
	}
}

func (lb *RateLimiter) couldTickerBeStoppedLocked() bool {
	return lb.availableInt == lb.capacity && lb.waitingCnt.Load() == 0
}
