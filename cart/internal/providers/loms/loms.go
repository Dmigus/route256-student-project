package loms

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/models"
)

type lomsClient interface {
	OrderCreate(ctx context.Context, userId int64, items []models.CartItem) (int64, error)
	GetNumberOfItemInStocks(ctx context.Context, skuId int64) (uint64, error)
}

type LOMS struct {
	client lomsClient
}

func NewLOMSProvider(client lomsClient) *LOMS {
	return &LOMS{
		client: client,
	}
}

// CreateOrder создаёт заказ для пользователя
func (L *LOMS) CreateOrder(ctx context.Context, userId int64, items []models.CartItem) (int64, error) {
	return L.client.OrderCreate(ctx, userId, items)
}

// IsItemAvailable проверяет, достаточно ли количество доступного для резервирования товара
func (L *LOMS) IsItemAvailable(ctx context.Context, skuId int64, count uint16) (bool, error) {
	available, err := L.client.GetNumberOfItemInStocks(ctx, skuId)
	if err != nil {
		return false, fmt.Errorf("could not get number of items with skuID %d in stocks: %w", skuId, err)
	}
	return available >= uint64(count), nil
}
