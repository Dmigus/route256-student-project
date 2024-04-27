// Package usecases представляет из себя реализацию сервисног слоя приложения
package loms

import (
	"context"
	"route256.ozon.ru/project/loms/internal/models"
)

type ordersCreator interface {
	Create(ctx context.Context, userId int64, items []models.OrderItem) (int64, error)
}

type ordersPayer interface {
	Pay(ctx context.Context, orderId int64) error
}

type stocksInfoGetter interface {
	GetNumOfAvailable(ctx context.Context, skuId int64) (uint64, error)
}

type ordersGetter interface {
	Get(ctx context.Context, orderId int64) (*models.Order, error)
}

type allOrdersGetter interface {
	Get(ctx context.Context) ([]*models.Order, error)
}

type ordersCanceller interface {
	Cancel(ctx context.Context, orderId int64) error
}

// LOMService представляет из себя объединение всех вариантов использования данного сервиса
type LOMService struct {
	ordersCreator    ordersCreator
	ordersPayer      ordersPayer
	stocksInfoGetter stocksInfoGetter
	ordersGetter     ordersGetter
	ordersCanceller  ordersCanceller
	allOrdersGetter  allOrdersGetter
}

// NewLOMService создаёт новый экземпляр LOMService
func NewLOMService(
	ordersCreator ordersCreator,
	ordersPayer ordersPayer,
	stocksInfoGetter stocksInfoGetter,
	ordersGetter ordersGetter,
	ordersCanceller ordersCanceller,
	allOrdersGetter allOrdersGetter) *LOMService {
	return &LOMService{
		ordersCreator:    ordersCreator,
		ordersPayer:      ordersPayer,
		stocksInfoGetter: stocksInfoGetter,
		ordersGetter:     ordersGetter,
		ordersCanceller:  ordersCanceller,
		allOrdersGetter:  allOrdersGetter,
	}
}

func (s *LOMService) CreateOrder(ctx context.Context, userId int64, items []models.OrderItem) (int64, error) {
	return s.ordersCreator.Create(ctx, userId, items)
}

func (s *LOMService) PayOrder(ctx context.Context, orderId int64) error {
	return s.ordersPayer.Pay(ctx, orderId)
}

func (s *LOMService) GetNumOfAvailable(ctx context.Context, skuId int64) (uint64, error) {
	return s.stocksInfoGetter.GetNumOfAvailable(ctx, skuId)
}

func (s *LOMService) GetOrder(ctx context.Context, orderId int64) (*models.Order, error) {
	return s.ordersGetter.Get(ctx, orderId)
}

func (s *LOMService) CancelOrder(ctx context.Context, orderId int64) error {
	return s.ordersCanceller.Cancel(ctx, orderId)
}

func (s *LOMService) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	return s.allOrdersGetter.Get(ctx)
}
