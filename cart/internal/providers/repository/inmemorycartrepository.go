package repository

import (
	"context"
	"go.opentelemetry.io/otel"
	"route256.ozon.ru/project/cart/internal/models"
	"sync"
)

var tracer = otel.Tracer("cart repository")

type InMemoryCartRepository struct {
	mu    sync.Mutex
	carts map[int64]*models.Cart
}

func New() *InMemoryCartRepository {
	return &InMemoryCartRepository{
		carts: make(map[int64]*models.Cart),
	}
}

func (i *InMemoryCartRepository) GetCart(ctx context.Context, user int64) (*models.Cart, error) {
	_, span := tracer.Start(ctx, "getting cart")
	defer span.End()
	i.mu.Lock()
	defer i.mu.Unlock()
	if cart, exists := i.carts[user]; exists {
		return cart, nil
	}
	span.AddEvent("creating new cart")
	newCart := models.NewCart()
	i.carts[user] = newCart
	return newCart, nil
}

func (i *InMemoryCartRepository) SaveCart(ctx context.Context, user int64, cart *models.Cart) error {
	_, span := tracer.Start(ctx, "saving cart")
	defer span.End()
	i.mu.Lock()
	defer i.mu.Unlock()
	i.carts[user] = cart
	return nil
}

// ClearCartReliable "надёжно" очищает корзину для пользователя. Сразу после возрата из метода корзина пользователя гарантированно пуста.
func (i *InMemoryCartRepository) ClearCartReliable(ctx context.Context, user int64) {
	ctx, span := tracer.Start(ctx, "clearing cart")
	defer span.End()
	cart, _ := i.GetCart(ctx, user)
	cart.Clear(ctx)
	i.mu.Lock()
	defer i.mu.Unlock()
	i.carts[user] = cart
}
