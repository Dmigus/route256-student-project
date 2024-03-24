package adder

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/cart/internal/models"
)

var errNotEnoughNumInStocks = errors.Wrap(models.ErrFailedPrecondition, "not enough item number in stocks")

type repository interface {
	GetCart(ctx context.Context, user int64) (*models.Cart, error)
	SaveCart(ctx context.Context, user int64, cart *models.Cart) error
}

type productService interface {
	IsItemPresent(ctx context.Context, skuId int64) (bool, error)
}

type stocksChecker interface {
	IsItemAvailable(ctx context.Context, skuId int64, count uint16) (bool, error)
}

// CartAdderService предназначен для добавления товаров в корзины пользователей
type CartAdderService struct {
	repo           repository
	productService productService
	stocks         stocksChecker
}

func New(repo repository, productService productService, stocks stocksChecker) *CartAdderService {
	return &CartAdderService{
		repo:           repo,
		productService: productService,
		stocks:         stocks,
	}
}

func (cs *CartAdderService) AddItem(ctx context.Context, user int64, skuId int64, count uint16) error {
	exists, err := cs.productService.IsItemPresent(ctx, skuId)
	if err != nil {
		return fmt.Errorf("could not check item %d presence: %w", skuId, err)
	}
	if !exists {
		return nil
	}
	isAvailable, err := cs.stocks.IsItemAvailable(ctx, skuId, count)
	if err != nil {
		return fmt.Errorf("could not check item %d availability: %w", skuId, err)
	}
	if !isAvailable {
		return errNotEnoughNumInStocks
	}
	cart, err := cs.repo.GetCart(ctx, user)
	if err != nil {
		return fmt.Errorf("could not get cart for user %d: %w", user, err)
	}
	cart.Add(ctx, skuId, count)
	err = cs.repo.SaveCart(ctx, user, cart)
	if err != nil {
		return fmt.Errorf("could not save cart for user %d: %w", user, err)
	}
	return nil
}
