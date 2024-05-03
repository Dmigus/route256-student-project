// Package ratelimiterhttp содержит round tripper с ограничением частоты запросов рейт лимитером
package ratelimiterhttp

import (
	"context"
	"net/http"
)

type (
	rateLimiter interface {
		Acquire(ctx context.Context) error
	}
	// RateLimitedTripper это структура, которая позволяет выполнять RoundTrip раунд трипера next с ограничением частоты запросов секунду
	RateLimitedTripper struct {
		rl   rateLimiter
		next http.RoundTripper
	}
)

// NewRateLimitedTripper создаёт новый экземпляр RateLimitedTripper
func NewRateLimitedTripper(rl rateLimiter, next http.RoundTripper) *RateLimitedTripper {
	return &RateLimitedTripper{
		rl:   rl,
		next: next,
	}
}

// RoundTrip обёртка над RoundTrip вложенного раунд трипера с ограничением по частоте. Если максимум достигунт, запрос будет блокирован до тех пор, пока
func (rlc *RateLimitedTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := rlc.rl.Acquire(r.Context()); err != nil {
		return nil, err
	}
	return rlc.next.RoundTrip(r)
}
