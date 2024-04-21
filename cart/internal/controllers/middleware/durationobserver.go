// Package middleware содержит функции, используемые в качестве http middleware
package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

const (
	// MethodNameLabel это метка, в которую запишется название http запроса
	MethodNameLabel = "method"
	// CodeLabel это метка в которую запишется код http ответа от cart
	CodeLabel = "code"
	// URLLabel это метка, содержашая очищенный url путь в запросе к cart
	URLLabel = "url"
)

type (
	observerVec interface {
		With(prometheus.Labels) prometheus.Observer
	}

	// DurationObserverMW это mw, который предназначен для подсчёта выполнения входящего запроса
	DurationObserverMW struct {
		wrapped     http.Handler
		observer    observerVec
		clearedPath string
	}
)

// NewDurationObserverMW создаёт DurationObserverMW
func NewDurationObserverMW(handlerToWrap http.Handler, observer observerVec, clearedPath string) *DurationObserverMW {
	return &DurationObserverMW{
		wrapped:     handlerToWrap,
		observer:    observer,
		clearedPath: clearedPath,
	}
}

// ServeHTTP обрабатывает запрос
func (do *DurationObserverMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := negroni.NewResponseWriter(w)
	start := time.Now()
	do.wrapped.ServeHTTP(lrw, r)
	end := time.Now()
	duration := end.Sub(start).Seconds()
	codeStr := strconv.Itoa(lrw.Status())
	labels := prometheus.Labels{MethodNameLabel: r.Method, CodeLabel: codeStr, URLLabel: do.clearedPath}
	do.observer.With(labels).Observe(duration)
}
