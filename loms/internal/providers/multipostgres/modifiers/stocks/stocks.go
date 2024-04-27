// Package stocks содержит реализацию стоков для транзакционной модификации данных в шардированном PostgreSQL.
package stocks

import (
	"context"

	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/stocks"
)

type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	// Stocks представляет репозиторий стоков с методами для модификации данных
	Stocks struct {
		trGetter     multipostgres.TransactionCreator
		shardManager shardmanager.Manager
		durRec       durationRecorder
	}
)

// NewStocks cоздаёт объект Stocks
func NewStocks(trGetter multipostgres.TransactionCreator, shardManager shardmanager.Manager, durRec durationRecorder) *Stocks {
	return &Stocks{trGetter: trGetter, shardManager: shardManager, durRec: durRec}
}

func (ps *Stocks) accessSinglePostgresImpl(ctx context.Context) (*stocks.Stocks, error) {
	shard := ps.shardManager.GetDefaultShard().Master()
	tr, err := ps.trGetter.GetTransaction(ctx, shard)
	if err != nil {
		return nil, err
	}
	return stocks.NewStocks(tr, ps.durRec), nil
}

// SetItemUnits устанавливает общее и зарезервированное количество товаров в стоках
func (ps *Stocks) SetItemUnits(ctx context.Context, skuID int64, total, reserved uint64) error {
	baseStocks, err := ps.accessSinglePostgresImpl(ctx)
	if err != nil {
		return err
	}
	return baseStocks.SetItemUnits(ctx, skuID, total, reserved)
}

// AddItems добавляет незарезервированные позиции к итемам.
func (ps *Stocks) AddItems(ctx context.Context, items []models.OrderItem) error {
	baseStocks, err := ps.accessSinglePostgresImpl(ctx)
	if err != nil {
		return err
	}
	return baseStocks.AddItems(ctx, items)
}

// Reserve резервирует товары. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *Stocks) Reserve(ctx context.Context, items []models.OrderItem) error {
	baseStocks, err := ps.accessSinglePostgresImpl(ctx)
	if err != nil {
		return err
	}
	return baseStocks.Reserve(ctx, items)
}

// RemoveReserved удаляет элементы из имеющихся и зарезервированных. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *Stocks) RemoveReserved(ctx context.Context, items []models.OrderItem) error {
	baseStocks, err := ps.accessSinglePostgresImpl(ctx)
	if err != nil {
		return err
	}
	return baseStocks.RemoveReserved(ctx, items)
}

// CancelReserved отменяет резервирование для итемов items. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *Stocks) CancelReserved(ctx context.Context, items []models.OrderItem) error {
	baseStocks, err := ps.accessSinglePostgresImpl(ctx)
	if err != nil {
		return err
	}
	return baseStocks.CancelReserved(ctx, items)
}
