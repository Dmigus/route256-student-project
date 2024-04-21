package mw

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Метки для секционирования результатов замеров времени
const (
	MethodNameLabel = "method"
	CodeLabel       = "code"
)

type (
	observerVec interface {
		With(prometheus.Labels) prometheus.Observer
	}
	// RequestDurationInterceptor это mw для замера и сохранения времени
	RequestDurationInterceptor struct {
		observer observerVec
	}
)

// NewRequestDurationInterceptor создаёт новый RequestDurationInterceptor, который будет сохранять результат в observer
func NewRequestDurationInterceptor(observer observerVec) *RequestDurationInterceptor {
	return &RequestDurationInterceptor{observer: observer}
}

// RecordDuration замеряет и сохраняет время запроса
func (r *RequestDurationInterceptor) RecordDuration(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	resp, err = handler(ctx, req)
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	code := status.Code(err)
	codeStr := strconv.FormatInt(int64(code), 10)
	labels := prometheus.Labels{MethodNameLabel: info.FullMethod, CodeLabel: codeStr}
	r.observer.With(labels).Observe(duration)
	return resp, err
}
