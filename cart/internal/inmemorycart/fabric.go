package inmemorycart

import (
	"context"
	repository "route256.ozon.ru/project/cart/internal/cartrepository"
)

type Fabric struct {
}

func (i *Fabric) Create(_ context.Context) (repository.Cart, error) {
	return New(), nil
}
