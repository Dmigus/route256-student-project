package middleware

import (
	"net/http"

	"go.uber.org/zap"
	"route256.ozon.ru/project/cart/internal/pkg/logging"
)

const (
	methodFieldName = "method"
	pathFieldName   = "path"
)

// RequestLoggerMW это mw для логирования информации о входящих http запросов в cart
type RequestLoggerMW struct {
	wrapped http.Handler
	logger  *zap.Logger
}

// NewLogger возвращает новый RequestLoggerMW
func NewLogger(handlerToWrap http.Handler, logger *zap.Logger) *RequestLoggerMW {
	return &RequestLoggerMW{
		wrapped: handlerToWrap,
		logger:  logger,
	}
}

// ServeHTTP обрабатывает запрос
func (rl *RequestLoggerMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fields := logging.AddTraceFieldsFromCtx(r.Context(),
		zap.String(methodFieldName, r.Method), zap.String(pathFieldName, r.URL.Path))
	rl.logger.Info("request received", fields...)
	rl.wrapped.ServeHTTP(w, r)
}
