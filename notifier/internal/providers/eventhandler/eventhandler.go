package eventhandler

import (
	"context"
	"io"
	"log"
	"route256.ozon.ru/project/notifier/internal/models"
)

type LoggerToWriter struct {
	logger *log.Logger
}

func NewLoggerToWriter(wr io.Writer) *LoggerToWriter {
	return &LoggerToWriter{
		logger: log.New(wr, "", log.Lmsgprefix|log.LstdFlags),
	}
}

func (s *LoggerToWriter) Handle(_ context.Context, event *models.OrderStatusChangedEvent) error {
	s.logger.Printf("Order with id = %d changed status: %s at %s\n", event.OrderID, event.Message, event.At)
	return nil
}
