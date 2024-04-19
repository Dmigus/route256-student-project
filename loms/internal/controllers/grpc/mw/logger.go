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

type LoggerMW struct {
	logger *zap.Logger
}

func NewLoggerMW(logger *zap.Logger) *LoggerMW {
	return &LoggerMW{logger: logger}
}

func (l *LoggerMW) LogReqAndResp(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if l.logger.Level().Enabled(zap.DebugLevel) {
		raw, _ := protojson.Marshal((req).(proto.Message))
		fields := logging.AddTraceFieldsFromCtx(ctx,
			zap.String(methodFieldName, info.FullMethod), zap.String("request", string(raw)))
		l.logger.Debug("request received", fields...)
	}
	if resp, err = handler(ctx, req); err != nil {
		fields := logging.AddTraceFieldsFromCtx(ctx,
			zap.String(methodFieldName, info.FullMethod), zap.Error(err))
		l.logger.Error("handler returned error", fields...)
		return
	}
	if l.logger.Level().Enabled(zap.DebugLevel) {
		rawResp, _ := protojson.Marshal((resp).(proto.Message))
		fields := logging.AddTraceFieldsFromCtx(ctx,
			zap.String(methodFieldName, info.FullMethod), zap.String("response", string(rawResp)))
		l.logger.Debug("response from handler", fields...)
	}
	return
}
