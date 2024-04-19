package logging

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func AddTraceFieldsFromCtx(ctx context.Context, fields ...zap.Field) []zap.Field {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		fields = append(fields, zap.String("trace_id", traceID.String()))
	}
	if spanCtx.HasSpanID() {
		spanID := spanCtx.SpanID()
		fields = append(fields, zap.String("span_id", spanID.String()))
	}
	return fields
}
