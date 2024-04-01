// Package stocksinfogetter содержит логику работы юзкейса получения информации об остатков товаров
package stocksinfogetter

import (
	"context"
	"fmt"
)

type (
	// StockRepo это контракт для использования репозитория стоков StocksInfoGetter'ом. Используется другими слоями для настройки доступа к исключительно зафиксированным данным
	StockRepo interface {
		GetNumOfAvailable(context.Context, int64) (uint64, error)
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, _ any, stocks StockRepo) error) error
	}
	// StocksInfoGetter - сущность, которая умеет извлекать данные об остатках
	StocksInfoGetter struct {
		tx txManager
	}
)

// NewGetter создаёт StocksInfoGetter. tx - должен быть объектом, позволяющим читать только зафиксированные данные
func NewGetter(tx txManager) *StocksInfoGetter {
	return &StocksInfoGetter{tx: tx}
}

// GetNumOfAvailable возвращает информацию об остатках товара с id = skuID
func (g *StocksInfoGetter) GetNumOfAvailable(ctx context.Context, skuID int64) (cnt uint64, err error) {
	err = g.tx.WithinTransaction(ctx, func(ctx context.Context, _ any, stocks StockRepo) error {
		cnt, err = stocks.GetNumOfAvailable(ctx, skuID)
		if err != nil {
			return fmt.Errorf("could not get number of available points for item %d: %w", skuID, err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
