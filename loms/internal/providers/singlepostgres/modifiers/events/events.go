// Package events содержит функциональность для работы с очередью сообщений, хранящихся в postgres
package events

import (
	"context"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/events/converter"

	"google.golang.org/protobuf/proto"
	"route256.ozon.ru/project/loms/internal/models"
)

const (
	messageOutboxTableName = "message_outbox"
)

// Events это структура, позволяющая работать с очередью событий
type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	Events struct {
		queries *Queries
		reqDur  durationRecorder
	}
)

// NewEvents cоздаёт объект объект очереди, работающего в рамках транзакции db
func NewEvents(db DBTX, reqDur durationRecorder) *Events {
	return &Events{queries: New(db), reqDur: reqDur}
}

// OrderStatusChanged сохраняет новое событие изменения статуса заказа
func (e *Events) OrderStatusChanged(ctx context.Context, order *models.Order) error {
	dt, err := e.currentDatetime(ctx)
	if err != nil {
		return err
	}
	eventMessage := converter.NewEventMessage(order, dt)
	partitionKey := []byte(strconv.FormatInt(order.Id(), 10))
	payload, err := proto.Marshal(eventMessage)
	if err != nil {
		return errors.Wrap(err, "could not marshal event message")
	}
	params := pushEventParams{PartitionKey: partitionKey, Payload: payload}
	e.reqDur.RecordDuration(messageOutboxTableName, sqlmetrics.Insert, func() error {
		err = e.queries.pushEvent(ctx, params)
		return err
	})
	return err
}

// PullNextEvents удаляет набор сообщений из outbox, не более  batchSize за раз и возвращает их в порядке
func (e *Events) PullNextEvents(ctx context.Context, batchSize int32) ([]models.EventMessage, error) {
	var events []MessageOutbox
	var err error
	e.reqDur.RecordDuration(messageOutboxTableName, sqlmetrics.Delete, func() error {
		events, err = e.queries.pullEvents(ctx, batchSize)
		return err
	})
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
