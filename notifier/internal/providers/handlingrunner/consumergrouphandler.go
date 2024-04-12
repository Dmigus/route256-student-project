package handlingrunner

import (
	"errors"
	"github.com/IBM/sarama"
	"route256.ozon.ru/project/notifier/internal/models"
	"route256.ozon.ru/project/notifier/internal/service"
	"strconv"
	"time"
)

var errMsgCorrupted = errors.New("message corrupted")

const atHeaderKey = "At"

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

func messageToEvent(message *sarama.ConsumerMessage) (*models.OrderStatusChangedEvent, error) {
	orderID, err := strconv.ParseInt(string(message.Key), 10, 64)
	if err != nil {
		return nil, errMsgCorrupted
	}
	headers := message.Headers
	if len(headers) == 0 {
		return nil, errMsgCorrupted
	}
	var at time.Time
	for _, header := range headers {
		if string(header.Key) == atHeaderKey {
			at, err = time.Parse(time.RFC3339, string(header.Value))
			break
		}
	}
	if err != nil {
		return nil, errMsgCorrupted
	}
	return &models.OrderStatusChangedEvent{
		OrderID: orderID,
		Message: string(message.Value),
		At:      at,
	}, nil
}
