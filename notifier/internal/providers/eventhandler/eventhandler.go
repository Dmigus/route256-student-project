// Package eventhandler содрежит обработчики событий
package eventhandler

import (
	"context"
	"io"
	"log"

	"route256.ozon.ru/project/notifier/internal/providers/eventhandler/converter"

	"route256.ozon.ru/project/notifier/internal/models"
)

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
func (s *LoggerToWriter) Handle(_ context.Context, message *models.EventMessage) error {
	evMess, err := converter.MessageToChangeOrderStatusEvent(message)
	if err != nil {
		return err
	}
	orderID := evMess.GetOrderID().GetOrderID()
	newStatus := converter.TransportStatusToString(evMess.GetStatus())
	dt := evMess.GetDatetime().AsTime()
	s.logger.Printf("Order with id = %d changed status to %s at %s\n", orderID, newStatus, dt)
	return nil
}
