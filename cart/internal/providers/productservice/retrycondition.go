package productservice

import "net/http"

// я не смог определить, куда следует положить логику определения условия ретрая для ProductService, поэтому положил сюда

const status420 = 420

func RetryCondition(resp *http.Response, err error) bool {
	if err != nil || resp == nil {
		return false
	}
	if resp.StatusCode == status420 || resp.StatusCode == http.StatusTooManyRequests {
		return true
	}
	return false
}
