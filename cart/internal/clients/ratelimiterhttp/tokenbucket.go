package ratelimiterhttp

import (
	"context"
	"time"
)

// TokenBucket это структура, которая позволяет ограничивать количество используемых ресурсов в единицу времени
type TokenBucket struct {
	available chan struct{}
}

// NewTokenBucket возвращает новый *TokenBucket, настроенный таким образом, что ресурс будет пополняться с частотой rps в секунду. Сразу после инициализации количество ресурсов полное.
func NewTokenBucket(ctx context.Context, rps int) *TokenBucket {
	avail := make(chan struct{}, rps)
	for i := 0; i < rps; i++ {
		avail <- struct{}{}
	}
	lb := &TokenBucket{available: avail}
	fillInterval := time.Second / time.Duration(rps)
	lb.runFiller(ctx, fillInterval)
	return lb
}

func (lb *TokenBucket) runFiller(ctx context.Context, fillInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(fillInterval)
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

func (lb *TokenBucket) addItemFree() {
	select {
	case lb.available <- struct{}{}:
	default:
	}
}

// Acquire это блокирующий вызов, который завершается в случаях:
// 1) если получилось использовать ресурс. В таком случае возращаемое значение nil
// 2) если завершился переданный контекст. В таком  случае возвращается прчина завершения контекста
func (lb *TokenBucket) Acquire(ctx context.Context) error {
	select {
	case <-lb.available:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
