package modifier

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

func TestAddItemWithoutErr(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		user  int64
		skuId int64
		count uint16
	}
	type testStruct struct {
		name string
		args args
	}
	tests := []testStruct{}
	mc := minimock.NewController(t)
	repo := NewRepositoryMock(mc)
	repo.SaveCartMock.Return(nil)
	prodServ := NewProductServiceMock(mc)

	// positive case
	arg := args{
		ctx:   context.Background(),
		user:  123,
		skuId: 123,
		count: 1,
	}
	cart := models.NewCart()
	repo.GetCartMock.When(minimock.AnyContext, arg.user).Then(cart, nil)
	prodServ.IsItemPresentMock.When(minimock.AnyContext, arg.skuId).Then(true, nil)
	tests = append(tests, testStruct{"positive", arg})

	// item not exists in product service
	arg = args{
		ctx:   context.Background(),
		user:  124,
		skuId: 124,
		count: 1,
	}
	prodServ.IsItemPresentMock.When(minimock.AnyContext, arg.skuId).Then(false, nil)
	tests = append(tests, testStruct{"item not exists in product service", arg})

	adder := New(repo, prodServ)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := adder.AddItem(tt.args.ctx, tt.args.user, tt.args.skuId, tt.args.count)
			assert.NoError(t, err, "must be no error")
		})
	}
}

func TestAddItemWithErrs(t *testing.T) {
	type args struct {
		ctx   context.Context
		user  int64
		skuId int64
		count uint16
	}

	type testStruct struct {
		name string
		args args
		err  error
	}
	tests := []testStruct{}
	mc := minimock.NewController(t)
	repo := NewRepositoryMock(mc)
	repo.SaveCartMock.Return(nil)
	prodServ := NewProductServiceMock(mc)

	//error checking good presence
	arg := args{
		ctx:   context.Background(),
		user:  123,
		skuId: 123,
		count: 1,
	}
	err := fmt.Errorf("someerror")
	prodServ.IsItemPresentMock.When(minimock.AnyContext, arg.skuId).Then(false, err)
	tests = append(tests, testStruct{"error checking good presence", arg, err})

	// error getting user cart
	arg = args{
		ctx:   context.Background(),
		user:  124,
		skuId: 124,
		count: 1,
	}
	err = fmt.Errorf("someerror2")
	prodServ.IsItemPresentMock.When(minimock.AnyContext, arg.skuId).Then(true, nil)
	repo.GetCartMock.When(minimock.AnyContext, arg.user).Then(nil, err)
	tests = append(tests, testStruct{"error getting user cart", arg, err})

	// error saving user cart
	arg = args{
		ctx:   context.Background(),
		user:  125,
		skuId: 125,
		count: 1,
	}
	err = fmt.Errorf("someerror3")
	cart := models.NewCart()
	prodServ.IsItemPresentMock.When(minimock.AnyContext, arg.skuId).Then(true, nil)
	repo.GetCartMock.When(minimock.AnyContext, arg.user).Then(cart, nil)
	repo.SaveCartMock.When(minimock.AnyContext, arg.user, cart).Then(err)
	tests = append(tests, testStruct{"error saving user cart", arg, err})

	adder := New(repo, prodServ)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := adder.AddItem(tt.args.ctx, tt.args.user, tt.args.skuId, tt.args.count)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
