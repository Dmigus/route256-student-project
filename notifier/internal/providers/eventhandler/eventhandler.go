// Package eventhandler содрежит обработчики событий
package eventhandler

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	v1 "route256.ozon.ru/project/notifier/internal/pkg/api/loms/v1"
	"route256.ozon.ru/project/notifier/internal/pkg/logging"

	"route256.ozon.ru/project/notifier/internal/providers/eventhandler/converter"

	"route256.ozon.ru/project/notifier/internal/models"
)

var tracer = otel.Tracer("event handler")

// LoggerToWriter это обработчик, который логирует поступившее событие в переданный Writer
type LoggerToWriter struct {
	logger *zap.Logger
}

// NewLoggerToWriter создаёт LoggerToWriter, который будет логировать в wr
func NewLoggerToWriter(logger *zap.Logger) *LoggerToWriter {
	return &LoggerToWriter{
		logger: logger,
	}
}

// Handle осуществляет обработку события
func (s *LoggerToWriter) Handle(ctx context.Context, message *models.EventMessage) error {
	ctx, span := tracer.Start(ctx, "eventMessage handling")
	defer span.End()
	evMess, err := converter.MessageToChangeOrderStatusEvent(message)
	if err != nil {
		span.SetStatus(codes.Error, "eventMessage was not recognized")
		return err
	}
	span.AddEvent("eventMessage recognized as ChangeOrderStatusEvent")
	s.handleChangeOrderStatusEvent(ctx, evMess)
	span.SetStatus(codes.Ok, "")
	return nil
}

func (s *LoggerToWriter) handleChangeOrderStatusEvent(ctx context.Context, ev *v1.ChangeOrderStatusEvent) {
	orderID := ev.GetOrderID().GetOrderID()
	userID := ev.GetUserID()
	newStatus := converter.TransportStatusToString(ev.GetStatus())
	dt := ev.GetDatetime().AsTime()
	fields := logging.AddTraceFieldsFromCtx(ctx,
		zap.Int64("order_id", orderID),
		zap.Int64("user_id", userID),
		zap.String("new_status", newStatus),
		zap.Time("at", dt),
	)
	s.logger.Info("Order status changed", fields...)
}
