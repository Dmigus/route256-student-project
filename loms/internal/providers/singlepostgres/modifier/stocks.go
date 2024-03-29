package modifier

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
	"sort"
)

var (
	errInsufficientStocks = errors.Wrap(models.ErrFailedPrecondition, "insufficient stocks")
	errItemIsNotFound     = errors.Wrap(models.ErrNotFound, "item is not found")
)

type Stocks struct {
	queries *Queries
}

func NewStocks(db DBTX) *Stocks {
	return &Stocks{queries: New(db)}
}

// SetItemUnits устанавливает общее и зарезервированное количество товаров в стоках
func (ps *Stocks) SetItemUnits(ctx context.Context, skuID int64, total, reserved uint64) error {
	params := insertStockParams{SkuID: skuID, Total: int32(total), Reserved: int32(reserved)}
	return ps.queries.insertStock(ctx, params)
}

// AddItems добавляет незарезервированные позиции к итемам.
func (ps *Stocks) AddItems(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		row, err := ps.queries.selectCount(ctx, it.SkuId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error adding %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newTotal := row.Total + int32(it.Count)
		params := updateTotalReservedParams{SkuID: it.SkuId, Total: newTotal, Reserved: row.Reserved}
		err = ps.queries.updateTotalReserved(ctx, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reserve резервирует товары. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *Stocks) Reserve(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		row, err := ps.queries.selectCount(ctx, it.SkuId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error adding %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newReserved := row.Reserved + int32(it.Count)
		if newReserved > row.Total {
			return fmt.Errorf("error reserving %d units if item with skuId = %d: %w", it.Count, it.SkuId, errInsufficientStocks)
		}
		params := updateReservedParams{SkuID: it.SkuId, Reserved: newReserved}
		err = ps.queries.updateReserved(ctx, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveReserved удаляет элементы из имеющихся и зарезервированных. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *Stocks) RemoveReserved(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		row, err := ps.queries.selectCount(ctx, it.SkuId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error adding %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newTotalCnt := row.Total - int32(it.Count)
		newReserved := row.Reserved - int32(it.Count)
		params := updateTotalReservedParams{SkuID: it.SkuId, Total: newTotalCnt, Reserved: newReserved}
		err = ps.queries.updateTotalReserved(ctx, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// CancelReserved отменяет резервирование для итемов items. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *Stocks) CancelReserved(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		row, err := ps.queries.selectCount(ctx, it.SkuId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error adding %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newReserved := row.Reserved - int32(it.Count)
		params := updateReservedParams{SkuID: it.SkuId, Reserved: newReserved}
		err = ps.queries.updateReserved(ctx, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSortedCopyItems(items []models.OrderItem) []models.OrderItem {
	sortedItems := make([]models.OrderItem, len(items))
	copy(sortedItems, items)
	sort.Slice(items, func(i, j int) bool {
		return items[i].SkuId < items[j].SkuId
	})
	return sortedItems
}
