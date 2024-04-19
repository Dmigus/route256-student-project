package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"route256.ozon.ru/project/cart/internal/pkg/logging"
)

const (
	methodFieldName = "method"
	pathFieldName   = "path"
)

type RequestLoggerMW struct {
	wrapped http.Handler
	logger  *zap.Logger
}

func NewLogger(handlerToWrap http.Handler, logger *zap.Logger) *RequestLoggerMW {
	return &RequestLoggerMW{
		wrapped: handlerToWrap,
		logger:  logger,
	}
}

func (rl *RequestLoggerMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fields := logging.AddTraceFieldsFromCtx(r.Context(),
		zap.String(methodFieldName, r.Method), zap.String(pathFieldName, r.URL.Path))
	rl.logger.Info("request received", fields...)
	rl.wrapped.ServeHTTP(w, r)
}
