package policies

import "net/http"

type RetryOnStatusCodes struct {
	retryCodes  []int
	maxAttempts int
}

func NewRetryOnStatusCodes(retryCodes []int, maxAttempts int) *RetryOnStatusCodes {
	return &RetryOnStatusCodes{
		retryCodes:  retryCodes,
		maxAttempts: maxAttempts,
	}
}

func (r *RetryOnStatusCodes) ShouldBeRetried(attemptNum int, _ *http.Request, resp *http.Response, respErr error) bool {
	if attemptNum >= r.maxAttempts {
		return false
	}
	if respErr != nil || resp == nil {
		return false
	}
	for _, badCode := range r.retryCodes {
		if badCode == resp.StatusCode {
			return true
		}
	}
	return false
}
