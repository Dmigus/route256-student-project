package in_memory_cart

import (
	"context"
	repository "route256.ozon.ru/project/cart/internal/cart-repository"
)

type InMemoryCartFabric struct {
}

func (i *InMemoryCartFabric) Create(ctx context.Context) (repository.Cart, error) {
	return New(), nil
}
