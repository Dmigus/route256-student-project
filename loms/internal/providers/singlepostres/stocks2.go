package singlepostres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

type PostgresStocks2 struct {
	tx pgx.Tx
}

func NewPostgresStocks2(tx pgx.Tx) *PostgresStocks2 {
	return &PostgresStocks2{tx: tx}
}

// SetItemUnits устанавливает общее и зарезервированное количество товаров в стоках
func (ps *PostgresStocks2) SetItemUnits(ctx context.Context, skuID int64, total, reserved uint64) error {
	_, err := ps.tx.Exec(ctx, insertStock, skuID, total, reserved)
	return err
}

// AddItems добавляет незарезервированные позиции к итемам.
func (ps *PostgresStocks2) AddItems(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := ps.tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error adding %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newTotal := totalCnt + uint64(it.Count)
		_, err = ps.tx.Exec(ctx, updateAll, it.SkuId, newTotal, reserved)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reserve резервирует товары. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *PostgresStocks2) Reserve(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := ps.tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
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
		_, err = ps.tx.Exec(ctx, updateReserved, it.SkuId, newReserved)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveReserved удаляет элементы из имеющихся и зарезервированных. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *PostgresStocks2) RemoveReserved(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := ps.tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error removing reservation %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newTotalCnt := totalCnt - uint64(it.Count)
		newReserved := reserved - uint64(it.Count)
		_, err = ps.tx.Exec(ctx, updateAll, it.SkuId, newTotalCnt, newReserved)
		if err != nil {
			return err
		}
	}
	return nil
}

// CancelReserved отменяет резервирование для итемов items. Если хотя бы один товар найти не удалось, то возращается ошибка с обозначением этого товара.
func (ps *PostgresStocks2) CancelReserved(ctx context.Context, items []models.OrderItem) error {
	sortedItems := getSortedCopyItems(items)
	for _, it := range sortedItems {
		var totalCnt, reserved uint64
		err := ps.tx.QueryRow(ctx, selectCountForUpdate, it.SkuId).Scan(&totalCnt, &reserved)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = errItemIsNotFound
			}
			return fmt.Errorf("error removing reservation %d units for item with skuId = %d: %w", it.Count, it.SkuId, err)
		}
		newReserved := reserved - uint64(it.Count)
		_, err = ps.tx.Exec(ctx, updateReserved, it.SkuId, newReserved)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetNumOfAvailable возвращает количество незарезервированных единиц для товара. Если не найден, будет ошибка.
func (ps *PostgresStocks2) GetNumOfAvailable(ctx context.Context, skuID int64) (uint64, error) {
	var totalCnt, reserved uint
	err := ps.tx.QueryRow(ctx, selectCount, skuID).Scan(&totalCnt, &reserved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errItemIsNotFound
		}
		return 0, err
	}
	return uint64(totalCnt - reserved), nil
}
