// Package durationobserverhttp содержит round tripper с наблидение времени запроса
package durationobserverhttp

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

const (
	MethodNameLabel = "method"
	CodeLabel       = "code"
	UrlLabel        = "url"
)

type (
	observerVec interface {
		With(prometheus.Labels) prometheus.Observer
	}
	// DurationObserverTripper это структура, которая позволяет выполнять RoundTrip раунд трипера next с ограничением частоты запросов секунду
	DurationObserverTripper struct {
		requestDurationObserver observerVec
		next                    http.RoundTripper
	}
)

// NewDurationObserverTripper создаёт новый экземпляр DurationObserverTripper
func NewDurationObserverTripper(reqDurationObserver observerVec, next http.RoundTripper) *DurationObserverTripper {
	return &DurationObserverTripper{
		requestDurationObserver: reqDurationObserver,
		next:                    next,
	}
}

// RoundTrip обёртка над RoundTrip вложенного раунд трипера с фиксированием длительности запроса
func (dot *DurationObserverTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	startTime := time.Now()
	resp, err := dot.next.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	codeStr := strconv.Itoa(resp.StatusCode)
	labels := prometheus.Labels{MethodNameLabel: r.Method, CodeLabel: codeStr, UrlLabel: r.URL.Path}
	dot.requestDurationObserver.With(labels).Observe(duration)
	return resp, nil
}
