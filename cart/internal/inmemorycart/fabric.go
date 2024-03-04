package inmemorycart

import (
	"context"
	repository "route256.ozon.ru/project/cart/internal/cartrepository"
)

type InMemoryCartFabric struct {
}

func (i *InMemoryCartFabric) Create(_ context.Context) (repository.Cart, error) {
	return New(), nil
}
