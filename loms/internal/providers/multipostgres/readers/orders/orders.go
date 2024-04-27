package orders

import (
	"context"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/readers/orders"
	"sort"
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

// LoadAll загружает информацию о всех заказах из всех шардов БД PostgreSQL
func (o *Orders) LoadAll(ctx context.Context) ([]*models.Order, error) {
	ordersByShards := make([][]*models.Order, o.shardManager.ShardNum())
	eg, ctx := errgroup.WithContext(ctx)
	for shardInd := 0; shardInd < o.shardManager.ShardNum(); shardInd++ {
		shard := o.shardManager.GetShardByInd(shardInd).Replica()
		eg.Go(func() error {
			shardOrders, err := o.loadAllForShard(ctx, shard)
			if err != nil {
				return err
			}
			ordersByShards[shardInd] = shardOrders
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	allOrders := lo.Flatten(ordersByShards)
	sort.Slice(allOrders, func(i, j int) bool {
		return allOrders[i].Id() > allOrders[j].Id()
	})
	return allOrders, nil
}

func (o *Orders) loadAllForShard(ctx context.Context, shard multipostgres.TxBeginner) ([]*models.Order, error) {
	tr, err := o.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return nil, err
	}
	baseOrders := orders.NewOrders(tr, o.reqDur)
	return baseOrders.LoadAll(ctx)
}
