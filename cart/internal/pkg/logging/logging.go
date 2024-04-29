// Package logging содержит вспомогательные функции, используемые для логирования в cart
package logging

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// AddTraceFieldsFromCtx дополняет поля fields полями trace_id и span_id из контекста ctx, если они там присутствуют
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

// EmptyLogger это пустой логер, который ничего не логирует
func EmptyLogger() *zap.Logger {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.FatalLevel),
		OutputPaths:      []string{},
		ErrorOutputPaths: []string{},
	}
	emptyL, _ := config.Build()
	return emptyL
}
