package usecases

import (
	"context"
	"errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var ErrService = errors.New("something went wrong during usecase execution")

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

type ordersCanceller interface {
	Cancel(ctx context.Context, orderId int64) error
}

type LOMService struct {
	ordersCreator    ordersCreator
	ordersPayer      ordersPayer
	stocksInfoGetter stocksInfoGetter
	ordersGetter     ordersGetter
	ordersCanceller  ordersCanceller
}

func NewLOMService(
	ordersCreator ordersCreator,
	ordersPayer ordersPayer,
	stocksInfoGetter stocksInfoGetter,
	ordersGetter ordersGetter,
	ordersCanceller ordersCanceller) *LOMService {
	return &LOMService{ordersCreator: ordersCreator, ordersPayer: ordersPayer, stocksInfoGetter: stocksInfoGetter, ordersGetter: ordersGetter, ordersCanceller: ordersCanceller}
}

func (s *LOMService) CreateOrder(ctx context.Context, userId int64, items []models.OrderItem) (int64, error) {
	res, err := s.ordersCreator.Create(ctx, userId, items)
	if err != nil {
		return 0, errors.Join(ErrService, err)
	}
	return res, nil
}

func (s *LOMService) PayOrder(ctx context.Context, orderId int64) error {
	if err := s.ordersPayer.Pay(ctx, orderId); err != nil {
		return errors.Join(ErrService, err)
	}
	return nil
}

func (s *LOMService) GetNumOfAvailable(ctx context.Context, skuId int64) (uint64, error) {
	res, err := s.stocksInfoGetter.GetNumOfAvailable(ctx, skuId)
	if err != nil {
		return 0, errors.Join(ErrService, err)
	}
	return res, nil
}

func (s *LOMService) GetOrder(ctx context.Context, orderId int64) (*models.Order, error) {
	res, err := s.ordersGetter.Get(ctx, orderId)
	if err != nil {
		return nil, errors.Join(ErrService, err)
	}
	return res, nil
}

func (s *LOMService) CancelOrder(ctx context.Context, orderId int64) error {
	if err := s.ordersCanceller.Cancel(ctx, orderId); err != nil {
		return errors.Join(ErrService, err)
	}
	return nil
}
