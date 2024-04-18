package mw

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
)

var tracer = otel.Tracer("mw")

func TraceRequest(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx, span := tracer.Start(ctx, "processing request by loms")
	defer span.End()
	resp, err := handler(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}
