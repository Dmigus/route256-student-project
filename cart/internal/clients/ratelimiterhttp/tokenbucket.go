package ratelimiterhttp

import (
	"context"
	"time"
)

// RateLimiter это структура, которая позволяет ограничивать количество используемых ресурсов в единицу времени
type RateLimiter struct {
	available    chan struct{}
	fillInterval time.Duration
}

// NewRateLimiter возвращает новый *RateLimiter, готовый к запуску.
// Сразу после инициализации количество ресурсов полное, то есть первые rps запросов Acquire будут удовлетворены без блокировки.
// После запуска ресурс будет пополняться с частотой rps в секунду.
func NewRateLimiter(rps int) *RateLimiter {
	avail := make(chan struct{}, rps)
	for i := 0; i < rps; i++ {
		avail <- struct{}{}
	}
	fillInterval := time.Second / time.Duration(rps)
	return &RateLimiter{available: avail, fillInterval: fillInterval}
}

func (lb *RateLimiter) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(lb.fillInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				lb.addItemFree()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (lb *RateLimiter) addItemFree() {
	select {
	case lb.available <- struct{}{}:
	default:
	}
}

// Acquire это блокирующий вызов, который завершается в случаях:
// 1) если получилось использовать ресурс. В таком случае возращаемое значение nil
// 2) если завершился переданный контекст. В таком  случае возвращается прчина завершения контекста
func (lb *RateLimiter) Acquire(ctx context.Context) error {
	select {
	case <-lb.available:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
