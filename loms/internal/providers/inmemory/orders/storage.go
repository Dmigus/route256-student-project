// Package orders содержит in-memory реализацию хранилища заказов
package orders

import (
	"context"
	"sync"

	pkgerrors "github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var errOrderNotFound = pkgerrors.Wrap(models.ErrNotFound, "order is not found")

type orderIDGenerator interface {
	NewID() int64
}

// InMemoryOrdersStorage представляет хранилище заказов
type InMemoryOrdersStorage struct {
	mu               sync.RWMutex
	data             map[int64]*models.Order
	orderIDGenerator orderIDGenerator
}

// NewInMemoryOrdersStorage создаёт новое хранилище InMemoryOrdersStorage. orderIDGenerator должен быть генератором уникальных id заказов
func NewInMemoryOrdersStorage(orderIDGenerator orderIDGenerator) *InMemoryOrdersStorage {
	return &InMemoryOrdersStorage{
		data:             make(map[int64]*models.Order),
		orderIDGenerator: orderIDGenerator,
	}
}

// Create создаёт заказ для юзера userID и товарами items в репозитории и возращает его
func (i *InMemoryOrdersStorage) Create(ctx context.Context, userID int64, items []models.OrderItem) (*models.Order, error) {
	newOrderID := i.orderIDGenerator.NewID()
	newOrder := models.NewOrder(userID, newOrderID)
	newOrder.Items = items
	err := i.Save(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}

// Save сохраняет заказ в хранилище
func (i *InMemoryOrdersStorage) Save(_ context.Context, order *models.Order) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data[order.Id()] = order
	return nil
}

// Load возвращает заказ по id. Если его нет, то errOrderNotFound
func (i *InMemoryOrdersStorage) Load(_ context.Context, orderID int64) (*models.Order, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	order, exists := i.data[orderID]
	if !exists {
		return nil, errOrderNotFound
	}
	return order, nil
}
