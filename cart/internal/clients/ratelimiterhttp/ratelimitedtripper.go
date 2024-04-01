package ratelimiterhttp

import (
	"go.uber.org/ratelimit"
	"net/http"
)

type RateLimitedTripper struct {
	rl   ratelimit.Limiter
	next http.RoundTripper
}

func NewRateLimitedTripper(rps int, next http.RoundTripper) *RateLimitedTripper {
	return &RateLimitedTripper{
		rl:   ratelimit.New(rps),
		next: next,
	}
}

func (rlc *RateLimitedTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rlc.rl.Take()
	return rlc.next.RoundTrip(r)
}
