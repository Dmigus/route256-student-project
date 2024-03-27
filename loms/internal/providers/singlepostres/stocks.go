package singlepostres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
	"sort"
)

var errInsufficientStocks = errors.Wrap(models.ErrFailedPrecondition, "insufficient stocks")
var errItemIsNotFound = errors.Wrap(models.ErrNotFound, "item is not found")

type PostgresStocks struct {
}

const (
	insertStock          = `INSERT INTO item_unit(sku_id, total, reserved) VALUES ($1, $2, $3) ON CONFLICT (sku_id) DO UPDATE SET total=$2, reserved=$3`
	selectCount          = `SELECT total, reserved FROM item_unit WHERE sku_id = $1`
	selectCountForUpdate = `SELECT total, reserved FROM item_unit WHERE sku_id = $1 FOR UPDATE`
	updateReserved       = `UPDATE item_unit SET reserved = $2 WHERE sku_id = $1`
	updateAll            = `UPDATE item_unit SET total = $2, reserved = $3 WHERE sku_id = $1`
)

func (ps *PostgresStocks) SetItemUnits(ctx context.Context, skuId int64, total, reserved uint64) error {
	tx := ctx.Value(trKey).(pgx.Tx)
	_, err := tx.Exec(ctx, insertStock, skuId, total, reserved)
	return err
}

func (ps *PostgresStocks) AddItems(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	tx := ctx.Value(trKey).(pgx.Tx)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error adding %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newTotal := totalCnt + uint64(it.Count)
		_, err = tx.Exec(ctx, updateAll, it.SkuId, newTotal, reserved)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PostgresStocks) Reserve(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	tx := ctx.Value(trKey).(pgx.Tx)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error reserving %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newReserved := reserved + uint64(it.Count)
		if newReserved > totalCnt {
			return fmt.Errorf("error reserving %d units if item with skuId = %d: %w", it.Count, it.SkuId, errInsufficientStocks)
		}
		_, err = tx.Exec(ctx, updateReserved, it.SkuId, newReserved)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PostgresStocks) RemoveReserved(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	tx := ctx.Value(trKey).(pgx.Tx)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error removing reservation %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newTotalCnt := totalCnt - uint64(it.Count)
		newReserved := reserved - uint64(it.Count)
		_, err = tx.Exec(ctx, updateAll, it.SkuId, newTotalCnt, newReserved)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PostgresStocks) CancelReserved(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	tx := ctx.Value(trKey).(pgx.Tx)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error removing reservation %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newReserved := reserved - uint64(it.Count)
		_, err = tx.Exec(ctx, updateReserved, it.SkuId, newReserved)
		if err != nil {
			return err
		}
	}
	return nil
}
func (ps *PostgresStocks) GetNumOfAvailable(ctx context.Context, skuId int64) (uint64, error) {
	tx := ctx.Value(trKey).(pgx.Tx)
	var totalCnt, reserved uint
	err := tx.QueryRow(ctx, selectCount, skuId).Scan(&totalCnt, &reserved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errItemIsNotFound
		}
		return 0, err
	}
	return uint64(totalCnt - reserved), nil
}

func getSortedCopyItems(items []models.OrderItem) []models.OrderItem {
	sortedItems := make([]models.OrderItem, len(items))
	copy(sortedItems, items)
	sort.Slice(items, func(i, j int) bool {
		return items[i].SkuId < items[j].SkuId
	})
	return sortedItems
}
