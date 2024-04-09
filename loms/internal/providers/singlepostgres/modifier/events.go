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
	params := insertEventParams{OrderID: order.Id(), Message: eventMessage}
	return e.queries.insertEvent(ctx, params)
}
