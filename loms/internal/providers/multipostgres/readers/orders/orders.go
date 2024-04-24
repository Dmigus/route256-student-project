package orders

import (
	"context"
	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/readers/orders"
)

// Orders представялет реализацию репозитория заказов с методами для чтения данных
type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	Orders struct {
		trGetter     multipostgres.TransactionCreator
		shardManager shardmanager.Manager
		reqDur       durationRecorder
	}
)

// NewOrders создаёт объект репозитория Orders
func NewOrders(trGetter multipostgres.TransactionCreator, shardManager shardmanager.Manager, reqDur durationRecorder) *Orders {
	return &Orders{trGetter: trGetter, shardManager: shardManager, reqDur: reqDur}
}

// Load загружает информацию о заказе из БД в PostgreSQL
func (o *Orders) Load(ctx context.Context, orderID int64) (*models.Order, error) {
	shardKey := multipostgres.OrderIDToShardKey(orderID)
	shard := o.shardManager.GetShard(shardKey).Replica()
	tr, err := o.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return nil, err
	}
	baseOrders := orders.NewOrders(tr, o.reqDur)
	return baseOrders.Load(ctx, orderID)
}
