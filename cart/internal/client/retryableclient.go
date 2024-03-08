package client

import (
	"context"
	"net/http"
)

type retryPolicy interface {
	ShouldBeRetried(attempts int, req *http.Request, resp *http.Response, respErr error) bool
}

type RetryableClient struct {
	client *http.Client
}

func NewRetryableClient(policy retryPolicy) *RetryableClient {
	retryRT := retryRoundTripper{
		next:   http.DefaultTransport,
		policy: policy,
	}
	return &RetryableClient{
		client: &http.Client{Transport: retryRT},
	}
}

func (rc *RetryableClient) Do(req *http.Request) (*http.Response, error) {
	return rc.client.Do(req)
}

type retryRoundTripper struct {
	next   http.RoundTripper
	policy retryPolicy
}

func (rr retryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
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
