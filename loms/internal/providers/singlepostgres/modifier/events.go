package modifier

import (
	"context"
	"fmt"

	"route256.ozon.ru/project/loms/internal/models"
)

const orderStatusChangedFmt = "order status was updated to %s"

// Events это структура, позволяющая работать с очередью событий
type Events struct {
	queries *Queries
}

// NewEvents cоздаёт объект объект очереди, работающего в рамках транзакции db
func NewEvents(db DBTX) *Events {
	return &Events{queries: New(db)}
}

// OrderStatusChanged сохраняет новое событие изменения статуса заказа
func (e *Events) OrderStatusChanged(ctx context.Context, order *models.Order) error {
	eventMessage := fmt.Sprintf(orderStatusChangedFmt, orderStatusToString(order.Status))
	params := pushEventParams{OrderID: order.Id(), Message: eventMessage}
	return e.queries.pushEvent(ctx, params)
}

// PullNextEvents удаляет набор событий из outbox, не более  batchSize за раз и возвращает их в порядке
func (e *Events) PullNextEvents(ctx context.Context, batchSize int32) ([]models.OrderStatusChangedEvent, error) {
	events, err := e.queries.pullEvents(ctx, batchSize)
	if err != nil {
		return nil, err
	}
	modelEvents := make([]models.OrderStatusChangedEvent, 0, len(events))
	for _, ev := range events {
		modelEvent := models.OrderStatusChangedEvent{OrderID: ev.OrderID, Message: ev.Message, At: ev.At.Time}
		modelEvents = append(modelEvents, modelEvent)
	}
	return modelEvents, nil
}
