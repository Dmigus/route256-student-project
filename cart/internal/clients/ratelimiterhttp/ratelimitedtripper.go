package ratelimiterhttp

import (
	"go.uber.org/ratelimit"
	"net/http"
)

// RateLimitedTripper это структура, которая позволяет выполнять RoundTrip раунд трипера next нен чаще rps запросов в секунду
type RateLimitedTripper struct {
	rl   ratelimit.Limiter
	next http.RoundTripper
}

// NewRateLimitedTripper создаёт новый экземпляр RateLimitedTripper
func NewRateLimitedTripper(rps int, next http.RoundTripper) *RateLimitedTripper {
	return &RateLimitedTripper{
		rl:   ratelimit.New(rps, ratelimit.WithSlack(0)),
		next: next,
	}
}

// RoundTrip обёртка над RoundTrip вложенного раунд трипера с ограничением по rps.
func (rlc *RateLimitedTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rlc.rl.Take()
	return rlc.next.RoundTrip(r)
}
