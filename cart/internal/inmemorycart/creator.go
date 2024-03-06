package inmemorycart

import (
	"context"
	repository "route256.ozon.ru/project/cart/internal/cartrepository"
)

type CartCreator struct {
}

func (i *CartCreator) Create(_ context.Context) (repository.Cart, error) {
	return New(), nil
}
