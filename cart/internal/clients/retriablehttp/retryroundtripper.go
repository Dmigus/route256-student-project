package retriablehttp

import (
	"context"
	"net/http"
)

type (
	retryPolicy interface {
		ShouldBeRetried(attempts int, req *http.Request, resp *http.Response, respErr error) bool
	}

	RetryRoundTripper struct {
		next   http.RoundTripper
		policy retryPolicy
	}
)

func NewRetryRoundTripper(next http.RoundTripper, policy retryPolicy) *RetryRoundTripper {
	return &RetryRoundTripper{
		next:   next,
		policy: policy,
	}
}

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
