// Package stocks содержит реализацию стоков только для чтения из PostgreSQL.
package stocks

import (
	"context"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var (
	errItemIsNotFound = errors.Wrap(models.ErrNotFound, "item is not found")
)

const itemUnitTable = "item_unit"

// Stocks представляет репозиторий стоков с методами для чтения данных
type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	Stocks struct {
		queries *Queries
		durRec  durationRecorder
	}
)

// NewStocks cоздаёт объект репозитория стоков, работающего в рамках транзакции tx
func NewStocks(db DBTX, durRec durationRecorder) *Stocks {
	return &Stocks{queries: New(db), durRec: durRec}
}

// GetNumOfAvailable возвращает количество незарезервированных единиц для товара. Если не найден, будет ошибка.
func (ps *Stocks) GetNumOfAvailable(ctx context.Context, skuID int64) (uint64, error) {
	var row selectCountRow
	var err error
	ps.durRec.RecordDuration(itemUnitTable, sqlmetrics.Select, func() error {
		row, err = ps.queries.selectCount(ctx, skuID)
		return err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errItemIsNotFound
		}
		return 0, err
	}
	return uint64(row.Total - row.Reserved), nil
}
