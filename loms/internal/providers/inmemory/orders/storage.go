package orders

import (
	"context"
	pkgerrors "github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
	"sync"
)

var errOrderNotFound = pkgerrors.Wrap(models.ErrNotFound, "order is not found")

type orderIDGenerator interface {
	NewID() int64
}

type InMemoryOrdersStorage struct {
	mu               sync.RWMutex
	data             map[int64]*models.Order
	orderIDGenerator orderIDGenerator
}

func NewInMemoryOrdersStorage(orderIDGenerator orderIDGenerator) *InMemoryOrdersStorage {
	return &InMemoryOrdersStorage{
		data:             make(map[int64]*models.Order),
		orderIDGenerator: orderIDGenerator,
	}
}

// Create создаёт заказ для юзера userID и товарами items в репозитории и возращает его
func (i *InMemoryOrdersStorage) Create(ctx context.Context, userID int64, items []models.OrderItem) (*models.Order, error) {
	newOrderId := i.orderIDGenerator.NewID()
	newOrder := models.NewOrder(userID, newOrderId)
	newOrder.Items = items
	err := i.Save(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (i *InMemoryOrdersStorage) Save(_ context.Context, order *models.Order) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data[order.Id()] = order
	return nil
}

func (i *InMemoryOrdersStorage) Load(_ context.Context, orderId int64) (*models.Order, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	order, exists := i.data[orderId]
	if !exists {
		return nil, errOrderNotFound
	}
	return order, nil
}
