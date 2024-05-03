// Package durationobserverhttp содержит round tripper с наблидение времени запроса
package durationobserverhttp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MethodNameLabel это метка в которую запишется название http запроса в product service
	MethodNameLabel = "method"
	// CodeLabel это метка в которую запишется код http ответа
	CodeLabel = "code"
	// URLLabel это метка, содержащая url путь в запросе к product service
	URLLabel = "url"
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
	labels := prometheus.Labels{MethodNameLabel: r.Method, CodeLabel: codeStr, URLLabel: r.URL.Path}
	dot.requestDurationObserver.With(labels).Observe(duration)
	return resp, nil
}
