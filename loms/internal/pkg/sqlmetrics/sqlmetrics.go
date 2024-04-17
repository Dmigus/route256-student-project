package sqlmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	CategoryLabel = "category"
	TableLabel    = "table"
	ErrLabel      = "code"
)

type (
	observerVec interface {
		With(prometheus.Labels) prometheus.Observer
	}
	SQLRequestDuration struct {
		observer observerVec
	}
)

func NewSQLRequestDuration(observer observerVec) *SQLRequestDuration {
	return &SQLRequestDuration{observer: observer}
}

// RecordDuration фиксирует длительность выполнения функции f в observer
func (r *SQLRequestDuration) RecordDuration(table string, category SQLCategory, f func() error) {
	startTime := time.Now()
	err := f()
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	labels := prometheus.Labels{TableLabel: table, CategoryLabel: category.String(), ErrLabel: errStr}
	r.observer.With(labels).Observe(duration)
}
