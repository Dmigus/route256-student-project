// Package sqlmetrics содержит функции для замера времени выполнения запросов в sql
package sqlmetrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// CategoryLabel это метка, в которой будет записана категория запроса
	CategoryLabel = "category"
	// TableLabel это метка, содержащая название таблицы, с которой производилась работа
	TableLabel = "table"
	// ErrLabel это метка, содержашая ошибку выполнения запроса
	ErrLabel = "error"
)

type (
	observerVec interface {
		With(prometheus.Labels) prometheus.Observer
	}
	// SQLRequestDuration это стуруктура, которая замеряет и записывает время выполнения запроса
	SQLRequestDuration struct {
		observer observerVec
	}
)

// NewSQLRequestDuration возвращает новый SQLRequestDuration, который будет записывать в observer
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
