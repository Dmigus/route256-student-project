// Package events содержит функциональность для работы с очередью сообщений, хранящихся в шардированном postgres
package events

import (
	"context"

	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/events"
)

type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	// EventsToInsert это структура, позволяющая вставлять в очередь событий
	EventsToInsert struct {
		trGetter     multipostgres.TransactionCreator
		shardManager shardmanager.Manager
		reqDur       durationRecorder
	}
	// EventsToPull это структура, позволяющая забирать из очереди событий
	EventsToPull struct {
		trGetter multipostgres.TransactionCreator
		shard    multipostgres.TxBeginner
		reqDur   durationRecorder
	}
)

// NewEventsToInsert cоздаёт объект EventsToInsert
func NewEventsToInsert(trGetter multipostgres.TransactionCreator, shardManager shardmanager.Manager, reqDur durationRecorder) *EventsToInsert {
	return &EventsToInsert{trGetter: trGetter, shardManager: shardManager, reqDur: reqDur}
}

// NewEventsToPull cоздаёт объект EventsToPull
func NewEventsToPull(trGetter multipostgres.TransactionCreator, shard multipostgres.TxBeginner, reqDur durationRecorder) *EventsToPull {
	return &EventsToPull{trGetter: trGetter, shard: shard, reqDur: reqDur}
}

// OrderStatusChanged сохраняет новое событие изменения статуса заказа
func (e *EventsToInsert) OrderStatusChanged(ctx context.Context, order *models.Order) error {
	shardKey := multipostgres.OrderIDToShardBucket(order.Id())
	shard := e.shardManager.GetShard(shardKey).Master()
	tr, err := e.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return err
	}
	baseEvents := events.NewEvents(tr, e.reqDur)
	return baseEvents.OrderStatusChanged(ctx, order)
}

// PullNextEvents удаляет набор самых первых сообщений из outbox, не более batchSize за раз и возвращает их
func (e *EventsToPull) PullNextEvents(ctx context.Context, batchSize int32) ([]models.EventMessage, error) {
	tr, err := e.trGetter.GetTransaction(ctx, e.shard)
	if err != nil {
		return nil, err
	}
	baseEvents := events.NewEvents(tr, e.reqDur)
	return baseEvents.PullNextEvents(ctx, batchSize)
}
