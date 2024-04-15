package handlingrunner

import (
	"github.com/IBM/sarama"
	"route256.ozon.ru/project/notifier/internal/models"
	"route256.ozon.ru/project/notifier/internal/service"
)

type consumerGroupHandler struct {
	handler service.EventHandler
}

func newConsumerGroupHandler(handler service.EventHandler) *consumerGroupHandler {
	return &consumerGroupHandler{handler: handler}
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
			// пока непонятно, куда логировать ошибки
			ev, _ := messageToEvent(message)
			err := c.handler.Handle(session.Context(), ev)
			if err == nil {
				session.MarkMessage(message, "")
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

func messageToEvent(message *sarama.ConsumerMessage) (*models.EventMessage, error) {
	return &models.EventMessage{
		PartitionKey: message.Key,
		Payload:      message.Value,
	}, nil
}
