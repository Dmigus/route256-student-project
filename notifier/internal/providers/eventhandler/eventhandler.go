// Package eventhandler содрежит обработчики событий
package eventhandler

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"io"
	"log"

	"route256.ozon.ru/project/notifier/internal/providers/eventhandler/converter"

	"route256.ozon.ru/project/notifier/internal/models"
)

var tracer = otel.Tracer("event handler")

// LoggerToWriter это обработчик, который логирует поступившее событие в переданный Writer
type LoggerToWriter struct {
	logger *log.Logger
}

// NewLoggerToWriter создаёт LoggerToWriter, который будет логировать в wr
func NewLoggerToWriter(wr io.Writer) *LoggerToWriter {
	prefix := "Handler is processing event: "
	return &LoggerToWriter{
		logger: log.New(wr, prefix, log.Lmsgprefix|log.LstdFlags),
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
	orderID := evMess.GetOrderID().GetOrderID()
	userID := evMess.GetUserID()
	newStatus := converter.TransportStatusToString(evMess.GetStatus())
	dt := evMess.GetDatetime().AsTime()
	s.logger.Printf("Order with id = %d (for user with id %d) changed status to %s at %s\n", orderID, userID, newStatus, dt)
	span.SetStatus(codes.Ok, "")
	return nil
}
