//go:build unit

package productinfogetter

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/cart/internal/models"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestProductInfoGetter_GetProductsInfoSuccess(t *testing.T) {
	// проверим, что после отработки теста нет работающих горутин
	defer goleak.VerifyNone(t)

	helper := newTestHelper(t)
	skuIDs := []int64{1, 2, 3}
	expected := []models.ProductInfo{
		{"first", 10},
		{"second", 20},
		{"third", 30},
	}
	helper.performMock.Set(func(_ context.Context, _ string, reqBody productservice.RequestWithSettableToken) (*GetProductResponse, error) {
		req := reqBody.(*GetProductRequest)
		var respBody GetProductResponse
		prodName := ""
		prodPrice := uint32(0)
		switch req.Sku {
		case 1:
			prodName = "first"
			prodPrice = 10
		case 2:
			prodName = "second"
			prodPrice = 20
		case 3:
			prodName = "third"
			prodPrice = 30
		}
		respBody.Price = &prodPrice
		respBody.Name = &prodName
		return &respBody, nil
	})
	returned, err := helper.prodInfoGetter.GetProductsInfo(context.Background(), skuIDs)
	require.NoError(t, err)
	assert.Equal(t, expected, returned)
}

func TestProductInfoGetter_GetProductsInfoError(t *testing.T) {
	helper := newTestHelper(t)
	skuIDs := []int64{1, 2}
	errorToThrow := fmt.Errorf("oops error")
	helper.performMock.Return(nil, errorToThrow)
	_, err := helper.prodInfoGetter.GetProductsInfo(context.Background(), skuIDs)
	assert.ErrorIs(t, err, errorToThrow)
}

// проверим, что GetProductsInfo инициирует набор запросов к PS с правильным контекстом.
func TestProductInfoGetter_GetProductsInfoContextCancellation(t *testing.T) {
	helper := newTestHelper(t)
	skuIDs := []int64{1, 2, 3}
	errorToThrow := fmt.Errorf("oops error")
	contextCancelTimeout := time.Second
	contextCancelled := &atomic.Bool{}
	wg := sync.WaitGroup{}
	wg.Add(len(skuIDs))
	helper.performMock.Set(func(ctx context.Context, _ string, reqBody productservice.RequestWithSettableToken) (*GetProductResponse, error) {
		defer wg.Done()
		req := reqBody.(*GetProductRequest)
		if req.Sku == 1 {
			return nil, errorToThrow
		}
		select {
		case <-ctx.Done():
			contextCancelled.Store(true)
			return nil, ctx.Err()
		case <-time.After(contextCancelTimeout):
		}
		var response GetProductResponse
		return &response, nil
	})
	_, err := helper.prodInfoGetter.GetProductsInfo(context.Background(), skuIDs)
	assert.ErrorIs(t, err, errorToThrow)
	// дождёмся выполнения всех горутин, чтобы получить результат о том, был ли отменён переданный им контекст
	wg.Wait()
	assert.True(t, contextCancelled.Load(), "context was not cancelled")
}
