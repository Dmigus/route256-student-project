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
	helper.performMock.Set(func(_ context.Context, _ string, reqBody productservice.RequestWithSettableToken, respBody any) (err error) {
		req := reqBody.(*getProductRequest)
		respToSet := respBody.(*getProductResponse)
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
		*respToSet = getProductResponse{Name: &prodName, Price: &prodPrice}
		return nil
	})
	returned, err := helper.prodInfoGetter.GetProductsInfo(context.Background(), skuIDs)
	require.NoError(t, err)
	assert.Equal(t, expected, returned)
}

func TestProductInfoGetter_GetProductsInfoError(t *testing.T) {
	helper := newTestHelper(t)
	skuIDs := []int64{1, 2}
	errorToThrow := fmt.Errorf("oops error")
	helper.performMock.Return(errorToThrow)
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
	helper.performMock.Set(func(ctx context.Context, _ string, reqBody productservice.RequestWithSettableToken, respBody any) (err error) {
		wg.Add(1)
		defer wg.Done()
		req := reqBody.(*getProductRequest)
		if req.Sku == 1 {
			return errorToThrow
		}
		select {
		case <-ctx.Done():
			contextCancelled.Store(true)
			return ctx.Err()
		case <-time.After(contextCancelTimeout):
		}
		respToSet := respBody.(*getProductResponse)
		prodName := ""
		prodPrice := uint32(0)
		respToSet.Price = &prodPrice
		respToSet.Name = &prodName
		return nil
	})
	_, err := helper.prodInfoGetter.GetProductsInfo(context.Background(), skuIDs)
	assert.ErrorIs(t, err, errorToThrow)
	// дождёмся выполнения всех горутин, чтобы получить результат о том, был ли отменён переданный им контекст
	wg.Wait()
	assert.True(t, contextCancelled.Load(), "context was not cancelled")
}
