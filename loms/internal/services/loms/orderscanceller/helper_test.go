//go:build unit
// +build unit

package orderscanceller

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/loms/internal/models"
	"sync/atomic"
	"testing"
)

type testHelper struct {
	stocksMock        *mStockRepoMockCancelReserved
	addItemsMock      *mStockRepoMockAddItems
	orderLoadRepoMock *mOrderRepoMockLoad
	orderSaveRepoMock *mOrderRepoMockSave
	canceller         *OrderCanceller
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrderRepoMock(mc)
	stocks := NewStockRepoMock(mc)
	evSender := NewEventSenderMock(mc)
	helper.orderLoadRepoMock = &(orders.LoadMock)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksMock = &(stocks.CancelReservedMock)
	helper.addItemsMock = &(stocks.AddItemsMock)
	trM := NewTxManagerMock(mc)
	trM.WithinTransactionMock.Set(func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo, stocks StockRepo, sender EventSender) bool) (err error) {
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
	helper.canceller = NewOrderCanceller(trM)
	return helper
}
