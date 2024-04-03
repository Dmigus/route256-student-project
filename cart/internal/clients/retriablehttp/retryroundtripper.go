// Package retriablehttp содержит round tripper с политикой повтора неуспешных запросов
package retriablehttp

import (
	"context"
	"net/http"
)

type (
	retryPolicy interface {
		ShouldBeRetried(attempts int, req *http.Request, resp *http.Response, respErr error) bool
	}
	// RetryRoundTripper это структура, которая позволяет выполнять RoundTrip раунд трипера next с политикой ретрая запроса
	RetryRoundTripper struct {
		next   http.RoundTripper
		policy retryPolicy
	}
)

// NewRetryRoundTripper создаёт новый экземпляр *RetryRoundTripper
func NewRetryRoundTripper(next http.RoundTripper, policy retryPolicy) *RetryRoundTripper {
	return &RetryRoundTripper{
		next:   next,
		policy: policy,
	}
}

// RoundTrip это обёртка над RoundTrip вложенного раунд трипера с проверкой ответа и повтором запроса, если это необходимо
func (rr RetryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error
	for attemptNum := 1; ; attemptNum++ {
		if contextWasDone(r.Context()) {
			return nil, r.Context().Err()
		}
		response, err = rr.next.RoundTrip(r)
		if !rr.policy.ShouldBeRetried(attemptNum, r, response, err) {
			return response, err
		}
	}
}

func contextWasDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
