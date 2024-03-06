package inmemorycart

import (
	"context"
	"route256.ozon.ru/project/cart/internal/providers/repository"
)

type CartCreator struct {
}

func (i *CartCreator) Create(_ context.Context) (repository.Cart, error) {
	return New(), nil
}
