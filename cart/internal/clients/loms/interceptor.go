package loms

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

const (
	MethodNameLabel = "method"
	CodeLabel       = "code"
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
