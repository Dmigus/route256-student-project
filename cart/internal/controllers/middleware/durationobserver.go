package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
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

	// DurationObserverMW это mw, который предназначен для подсчёта выполнения входящего запроса
	DurationObserverMW struct {
		wrapped     http.Handler
		observer    observerVec
		clearedPath string
	}
)

func NewDurationObserverMW(handlerToWrap http.Handler, observer observerVec, clearedPath string) *DurationObserverMW {
	return &DurationObserverMW{
		wrapped:     handlerToWrap,
		observer:    observer,
		clearedPath: clearedPath,
	}
}

func (do *DurationObserverMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := negroni.NewResponseWriter(w)
	start := time.Now()
	do.wrapped.ServeHTTP(lrw, r)
	end := time.Now()
	duration := end.Sub(start).Seconds()
	codeStr := strconv.Itoa(lrw.Status())
	labels := prometheus.Labels{MethodNameLabel: r.Method, CodeLabel: codeStr, UrlLabel: do.clearedPath}
	do.observer.With(labels).Observe(duration)
}
