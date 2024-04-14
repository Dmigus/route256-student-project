package events

import (
	"context"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifier/events/converter"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
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
	dt, err := e.currentDatetime(ctx)
	eventMessage := converter.NewEventMessage(order, dt)
	partitionKey := []byte(strconv.FormatInt(order.Id(), 10))
	payload, err := proto.Marshal(eventMessage)
	if err != nil {
		return errors.Wrap(err, "could not marshal event message")
	}
	params := pushEventParams{PartitionKey: partitionKey, Payload: payload}
	return e.queries.pushEvent(ctx, params)
}

// PullNextEvents удаляет набор сообщений из outbox, не более  batchSize за раз и возвращает их в порядке
func (e *Events) PullNextEvents(ctx context.Context, batchSize int32) ([]models.EventMessage, error) {
	events, err := e.queries.pullEvents(ctx, batchSize)
	if err != nil {
		return nil, err
	}
	modelEvents := make([]models.EventMessage, 0, len(events))
	for _, ev := range events {
		modelEvent := models.EventMessage{PartitionKey: ev.PartitionKey, Payload: ev.Payload}
		modelEvents = append(modelEvents, modelEvent)
	}
	return modelEvents, nil
}

func (e *Events) currentDatetime(ctx context.Context) (time.Time, error) {
	pgxDt, err := e.queries.currentDatetime(ctx)
	if err != nil || !pgxDt.Valid {
		return time.Time{}, errors.Wrap(err, "could not get current datetime")
	}
	return pgxDt.Time, nil
}
