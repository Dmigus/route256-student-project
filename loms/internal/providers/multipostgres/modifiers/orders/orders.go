// Package orders содержит реализацию заказов для транзакционной модификации данных в шардированном PostgreSQL.
package orders

import (
	"context"
	"math/rand"

	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/orders"
)

const (
	orderTableName     = "order"
	orderItemTableName = "order_item"
)

type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	// Orders представялет реализацию репозитория заказов с методами для модификации данных
	Orders struct {
		trGetter     multipostgres.TransactionCreator
		shardManager shardmanager.Manager
		durRec       durationRecorder
	}
)

// NewOrders создаёт объект репозитория заказов Orders
func NewOrders(trGetter multipostgres.TransactionCreator, shardManager shardmanager.Manager, durRec durationRecorder) *Orders {
	return &Orders{trGetter: trGetter, shardManager: shardManager, durRec: durRec}
}

// Create создаёт заказ для юзера userID и товарами items в репозитории и возращает его
func (po *Orders) Create(ctx context.Context, userID int64, items []models.OrderItem) (*models.Order, error) {
	shardKey := po.chooseShardKeyToNewOrder(userID)
	var err error
	queries, err := po.initQueriesForNewOrder(ctx, shardKey)
	if err != nil {
		return nil, err
	}
	var orderID int64
	params := createOrderParams{Column1: int(shardKey), UserID: userID, Status: "New", AreItemsReserved: false}
	po.durRec.RecordDuration(orderTableName, sqlmetrics.Insert, func() error {
		orderID, err = queries.createOrder(ctx, params)
		return err
	})
	if err != nil {
		return nil, err
	}
	order := models.NewOrder(userID, orderID)
	itemsParams := insertItemParamsFrom(orderID, items)
	po.durRec.RecordDuration(orderItemTableName, sqlmetrics.Insert, func() error {
		_, err = queries.insertOrderItem(ctx, itemsParams)
		return err
	})
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (po *Orders) initQueriesForNewOrder(ctx context.Context, shardKey shardmanager.ShardBucket) (*Queries, error) {
	shard := po.shardManager.GetShard(shardKey).Master()
	tx, err := po.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return nil, err
	}
	return New(tx), nil
}

func insertItemParamsFrom(orderID int64, items []models.OrderItem) []insertOrderItemParams {
	itemsParams := make([]insertOrderItemParams, 0, len(items))
	for _, it := range items {
		params := insertOrderItemParams{OrderID: orderID, SkuID: it.SkuId, Count: int32(it.Count)}
		itemsParams = append(itemsParams, params)
	}
	return itemsParams
}

func (po *Orders) chooseShardKeyToNewOrder(_ int64) shardmanager.ShardBucket {
	bucket := rand.Intn(shardmanager.BucketsNum)
	return shardmanager.ShardBucket(bucket)
}

// Save сохраняет заказ в БД в PostgreSQL. Изменение позиций заказа не предусмотрено
func (po *Orders) Save(ctx context.Context, order *models.Order) error {
	shardKey := multipostgres.OrderIDToShardBucket(order.Id())
	shard := po.shardManager.GetShard(shardKey).Master()
	tr, err := po.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return err
	}
	baseOrders := orders.NewOrders(tr, po.durRec)
	return baseOrders.Save(ctx, order)
}

// Load загружает информацию о заказе из БД в PostgreSQL, производя SELECT FOR UPDATE
func (po *Orders) Load(ctx context.Context, orderID int64) (*models.Order, error) {
	shardKey := multipostgres.OrderIDToShardBucket(orderID)
	shard := po.shardManager.GetShard(shardKey).Master()
	tr, err := po.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return nil, err
	}
	baseOrders := orders.NewOrders(tr, po.durRec)
	return baseOrders.Load(ctx, orderID)
}
