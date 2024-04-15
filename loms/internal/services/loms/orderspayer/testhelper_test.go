//go:build unit
// +build unit

package orderspayer

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/loms/internal/models"
	"sync/atomic"
	"testing"
)

type testHelper struct {
	orderLoadRepoMock *mOrderRepoMockLoad
	orderSaveRepoMock *mOrderRepoMockSave
	stocksRepoMock    *mStockRepoMockRemoveReserved
	payer             *OrdersPayer
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrderRepoMock(mc)
	stocks := NewStockRepoMock(mc)
	evSender := NewEventSenderMock(mc)
	helper.orderLoadRepoMock = &(orders.LoadMock)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksRepoMock = &(stocks.RemoveReservedMock)
	txM := NewTxManagerMock(t)
	txM.WithinTransactionMock.Set(func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo, stocks StockRepo, sender EventSender) bool) error {
		// выполнение функции
		f1(ctx, orders, stocks, evSender)
		// ошибки фиксации транзакции тестировать пока не будем
		return nil
	})
	// evSender пока всегда сохраняет без ошибок
	evSender.OrderStatusChangedMock.Set(func(ctx context.Context, _ *models.Order) (err error) {
		return nil
	})
	atomic.CompareAndSwapUint64(&evSender.afterOrderStatusChangedCounter, 0, 1)
	helper.payer = NewOrdersPayer(txM)
	return helper
}
