//go:build unit

package orderscreator

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/loms/internal/models"
	"testing"
)

func TestOrdersCreatorPositive(t *testing.T) {
	t.Parallel()
	h := newTestHelper(t)
	h.generatorMock.Expect().Return(1)
	items := []models.OrderItem{{123, 5}}
	h.stocksMock.Expect(minimock.AnyContext, items).Return(nil)
	h.orderSaveRepoMock.Return(nil)
	orderId, err := h.creator.Create(context.Background(), 123, items)
	require.NoError(t, err)
	assert.Equal(t, int64(1), orderId)
}

func TestOrdersCreatorErrors(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		userId int64
		items  []models.OrderItem
	}
	errorToThrow := fmt.Errorf("oops error")
	tests := []struct {
		name      string
		mockSetup func(testHelper)
		args      args
		err       error
	}{
		{
			name: "error reserving items",
			mockSetup: func(helper testHelper) {
				items := []models.OrderItem{{123, 5}}
				helper.generatorMock.Expect().Return(1)
				helper.stocksMock.Expect(minimock.AnyContext, items).Return(errorToThrow)
				helper.orderSaveRepoMock.Return(nil)
			},
			args: args{
				ctx:    context.Background(),
				userId: 123,
				items:  []models.OrderItem{{123, 5}},
			},
			err: errorToThrow,
		},
		{
			name: "error saving order",
			mockSetup: func(helper testHelper) {
				items := []models.OrderItem{{123, 5}}
				helper.generatorMock.Expect().Return(1)
				helper.stocksMock.Expect(minimock.AnyContext, items).Return(nil)
				helper.orderSaveRepoMock.Return(errorToThrow)
			},
			args: args{
				ctx:    context.Background(),
				userId: 123,
				items:  []models.OrderItem{{123, 5}},
			},
			err: errorToThrow,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := newTestHelper(t)
			tt.mockSetup(helper)
			_, err := helper.creator.Create(tt.args.ctx, tt.args.userId, tt.args.items)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestOrdersCreator_createOrderInstance(t *testing.T) {
	t.Parallel()
	h := newTestHelper(t)
	h.generatorMock.Expect().Return(15)
	newOrder := h.creator.createOrderInstance(1234)
	assert.Equal(t, int64(15), newOrder.Id())
	assert.Equal(t, models.New, newOrder.Status)
}
