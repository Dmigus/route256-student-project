// Package mw содержит middleware для grpc сервера
package mw

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"route256.ozon.ru/project/loms/internal/pkg/logging"
)

const methodFieldName = "method"

// LoggerMW это mw, предназначенынй для логирования входящих запросов и ответов от сервиса loms
type LoggerMW struct {
	logger *zap.Logger
}

// NewLoggerMW создаёт новый LoggerMW, который будет использовать logger
func NewLoggerMW(logger *zap.Logger) *LoggerMW {
	return &LoggerMW{logger: logger}
}

// LogReqAndResp производит логирование входящих запросов и ответов от сервиса loms
func (l *LoggerMW) LogReqAndResp(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if l.logger.Level().Enabled(zap.DebugLevel) {
		l.logReq(ctx, req, info.FullMethod)
	}
	if resp, err = handler(ctx, req); err != nil {
		l.logErr(ctx, err, info.FullMethod)
		return
	}
	if l.logger.Level().Enabled(zap.DebugLevel) {
		l.logResp(ctx, resp, info.FullMethod)
	}
	return
}

func (l *LoggerMW) logReq(ctx context.Context, req interface{}, methodName string) {
	raw, _ := protojson.Marshal((req).(proto.Message))
	fields := logging.AddTraceFieldsFromCtx(ctx,
		zap.String(methodFieldName, methodName), zap.String("request", string(raw)))
	l.logger.Debug("request received", fields...)
}

func (l *LoggerMW) logResp(ctx context.Context, resp interface{}, methodName string) {
	rawResp, _ := protojson.Marshal((resp).(proto.Message))
	fields := logging.AddTraceFieldsFromCtx(ctx,
		zap.String(methodFieldName, methodName), zap.String("response", string(rawResp)))
	l.logger.Debug("response from handler", fields...)
}

func (l *LoggerMW) logErr(ctx context.Context, err error, methodName string) {
	fields := logging.AddTraceFieldsFromCtx(ctx,
		zap.String(methodFieldName, methodName), zap.Error(err))
	l.logger.Error("handler returned error", fields...)
}
