// Package client содержит реализацию grpc клиента к сервису loms
package client

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/cart/internal/models"
	"route256.ozon.ru/project/cart/internal/providers/loms/client/converter"
	v1 "route256.ozon.ru/project/cart/internal/providers/loms/client/protoc/v1"
)

var (
	errNotFound           = errors.Wrap(models.ErrNotFound, "LOMS returned NotFound code")
	errFailedPrecondition = errors.Wrap(models.ErrFailedPrecondition, "LOMS returned FailedPrecondition code")
)

type (
	observerVec interface {
		With(prometheus.Labels) prometheus.Observer
	}
	// Client это структура для работы с loms
	Client struct {
		client v1.LOMServiceClient
	}
)

// NewClient возвращает новый Client, который, помимо прочего, записывает продолжительность запросов в метрику reqDurationObserver
func NewClient(addr string, reqDurationObserver observerVec) (*Client, error) {
	interceptor := requestDurationInterceptor{reqDurationObserver}
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.intercept),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	conn.Connect()
	if err != nil {
		return nil, err
	}
	return &Client{
		client: v1.NewLOMServiceClient(conn),
	}, nil
}

// OrderCreate создаёт новый заказ дял пользователя userId с итемами items
func (c *Client) OrderCreate(ctx context.Context, userId int64, items []models.CartItem) (int64, error) {
	request := converter.ModelsToOrderCreateRequest(userId, items)
	response, err := c.client.OrderCreate(ctx, request)
	if err != nil {
		err = detectKnownErrors(err)
		return 0, fmt.Errorf("error calling OrderCreate for user %d: %w", userId, err)
	}
	return converter.OrderIdToId(response), nil
}

// GetNumberOfItemInStocks возвращает количество товаров с id = skuId в стоках
func (c *Client) GetNumberOfItemInStocks(ctx context.Context, skuId int64) (uint64, error) {
	req := converter.SkuIdToStocksInfoRequest(skuId)
	response, err := c.client.StocksInfo(ctx, req)
	if err != nil {
		err = detectKnownErrors(err)
		return 0, fmt.Errorf("error calling StocksInfo for item %d: %w", skuId, err)
	}
	return converter.StocksInfoResponseToCount(response), nil
}

func detectKnownErrors(errResp error) error {
	code := status.Code(errResp)
	switch code {
	case codes.NotFound:
		return errors.Wrap(errNotFound, errResp.Error())
	case codes.FailedPrecondition:
		return errors.Wrap(errFailedPrecondition, errResp.Error())
	}
	return errResp
}
