package handlingrunner

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"route256.ozon.ru/project/notifier/internal/models"
	"route256.ozon.ru/project/notifier/internal/service"
)

type consumerGroupHandler struct {
	handler service.EventHandler
	logger  *zap.Logger
}

func newConsumerGroupHandler(handler service.EventHandler, logger *zap.Logger) *consumerGroupHandler {
	return &consumerGroupHandler{handler: handler, logger: logger}
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает до тех пор пока сессия не завершилась
func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			err := c.processTracedMessage(session.Context(), message)
			if err != nil {
				c.logger.Error("error handling message", zap.Error(err))
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *consumerGroupHandler) processTracedMessage(sessionCtx context.Context, message *sarama.ConsumerMessage) error {
	ctx := otel.GetTextMapPropagator().Extract(sessionCtx, otelsarama.NewConsumerMessageCarrier(message))
	ev := messageToEvent(message)
	return c.handler.Handle(ctx, ev)
}

func messageToEvent(message *sarama.ConsumerMessage) *models.EventMessage {
	return &models.EventMessage{
		PartitionKey: message.Key,
		Payload:      message.Value,
	}
}
