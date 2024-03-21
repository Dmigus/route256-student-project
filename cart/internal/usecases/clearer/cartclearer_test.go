//go:build unit

package clearer

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

func TestCartModifierService_ClearCart(t *testing.T) {
	errorToThrow := fmt.Errorf("oops error")
	type args struct {
		ctx  context.Context
		user int64
	}
	tests := []struct {
		name      string
		mockSetup func(testHelper)
		args      args
		wantErr   error
	}{
		{
			name: "positive",
			mockSetup: func(helper testHelper) {
				helper.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(models.NewCart(), nil)
				helper.saveCartRepoMock.Return(nil)
			},
			args: args{
				context.Background(),
				123,
			},
			wantErr: nil,
		},
		{
			name: "error getting user cart",
			mockSetup: func(helper testHelper) {
				helper.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(nil, errorToThrow)
			},
			args: args{
				context.Background(),
				123,
			},
			wantErr: errorToThrow,
		},
		{
			name: "error saving user cart",
			mockSetup: func(helper testHelper) {
				helper.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(models.NewCart(), nil)
				helper.saveCartRepoMock.Return(errorToThrow)
			},
			args: args{
				context.Background(),
				123,
			},
			wantErr: errorToThrow,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := newTestHelper(t)
			tt.mockSetup(helper)
			err := helper.service.ClearCart(tt.args.ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
