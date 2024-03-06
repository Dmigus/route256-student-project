package retryableclient

import "net/http"

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
		select {
		// проверка, что контекст к этому моменту не отменён
		case <-r.Context().Done():
			return nil, r.Context().Err()
		default:
			res, err = rr.next.RoundTrip(r)
			if !rr.retryCondition(res, err) {
				return res, err
			}
		}
	}
	return res, err
}
