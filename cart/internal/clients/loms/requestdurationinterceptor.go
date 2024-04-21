package loms

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	// MethodNameLabel это метка в которую запишется название метода loms
	MethodNameLabel = "method"
	// CodeLabel это метка в которую запишется код grpc ответа
	CodeLabel = "code"
)

type requestDurationInterceptor struct {
	requestDurationObserver observerVec
}

func (r *requestDurationInterceptor) intercept(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	startTime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	code := status.Code(err)
	codeStr := strconv.FormatInt(int64(code), 10)
	labels := prometheus.Labels{MethodNameLabel: method, CodeLabel: codeStr}
	r.requestDurationObserver.With(labels).Observe(duration)
	return err
}
