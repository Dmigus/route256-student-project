package stocks

import (
	"context"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/readers/stocks"
)

type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	// Stocks представляет репозиторий стоков с методами для чтения данных в рамках распределённой транзакции
	Stocks struct {
		trGetter     multipostgres.TransactionCreator
		shardManager shardmanager.Manager
		durRec       durationRecorder
	}
)

// NewStocks cоздаёт объект репозитория стоков, работающего в рамках транзакции tx
func NewStocks(trGetter multipostgres.TransactionCreator, shardManager shardmanager.Manager, durRec durationRecorder) *Stocks {
	return &Stocks{trGetter: trGetter, shardManager: shardManager, durRec: durRec}
}

// GetNumOfAvailable возвращает количество незарезервированных единиц для товара. Если не найден, будет ошибка.
func (s *Stocks) GetNumOfAvailable(ctx context.Context, skuID int64) (uint64, error) {
	shard := s.shardManager.GetDefaultShard().Replica()
	tr, err := s.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return 0, err
	}
	baseStocks := stocks.NewStocks(tr, s.durRec)
	return baseStocks.GetNumOfAvailable(ctx, skuID)
}
