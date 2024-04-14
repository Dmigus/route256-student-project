// Package stocks содержит реализацию стоков только для чтения из PostgreSQL.
package stocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var (
	errItemIsNotFound = errors.Wrap(models.ErrNotFound, "item is not found")
)

// Stocks представляет репозиторий стоков с методами для чтения данных
type Stocks struct {
	queries *Queries
}

// NewStocks cоздаёт объект репозитория стоков, работающего в рамках транзакции tx
func NewStocks(db DBTX) *Stocks {
	return &Stocks{queries: New(db)}
}

// GetNumOfAvailable возвращает количество незарезервированных единиц для товара. Если не найден, будет ошибка.
func (ps *Stocks) GetNumOfAvailable(ctx context.Context, skuID int64) (uint64, error) {
	row, err := ps.queries.selectCount(ctx, skuID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errItemIsNotFound
		}
		return 0, err
	}
	return uint64(row.Total - row.Reserved), nil
}
