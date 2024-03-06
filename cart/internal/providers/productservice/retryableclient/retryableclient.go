package retryableclient

import (
	"context"
	"net/http"
)

type RetryableClient struct {
	client *http.Client
}

func NewRetryableClient(maxRetries int, retryCondition func(*http.Response, error) bool) *RetryableClient {
	retryRT := retryRoundTripper{
		next:           http.DefaultTransport,
		maxRetries:     maxRetries,
		retryCondition: retryCondition,
	}
	return &RetryableClient{
		client: &http.Client{Transport: retryRT},
	}
}

func (rc *RetryableClient) Do(req *http.Request) (*http.Response, error) {
	return rc.client.Do(req)
}

type retryRoundTripper struct {
	next           http.RoundTripper
	maxRetries     int
	retryCondition func(*http.Response, error) bool
}

func (rr retryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error
	for attempts := 0; attempts < rr.maxRetries; attempts++ {
		if contextWasDone(r.Context()) {
			return nil, r.Context().Err()
		}
		res, err = rr.next.RoundTrip(r)
		if !rr.retryCondition(res, err) {
			return res, err
		}
	}
	return res, err
}

func contextWasDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
