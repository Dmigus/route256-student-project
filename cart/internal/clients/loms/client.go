package loms

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/cart/internal/clients/loms/converter"
	v1 "route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"
	"route256.ozon.ru/project/cart/internal/models"
)

var (
	errNotFound           = errors.Wrap(models.ErrNotFound, "LOMS returned NotFound code")
	errFailedPrecondition = errors.Wrap(models.ErrFailedPrecondition, "LOMS returned FailedPrecondition code")
)

type Client struct {
	client v1.LOMServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn.Connect()
	if err != nil {
		return nil, err
	}
	return &Client{
		client: v1.NewLOMServiceClient(conn),
	}, nil
}

func (c *Client) OrderCreate(ctx context.Context, userId int64, items []models.CartItem) (int64, error) {
	request := converter.ModelsToOrderCreateRequest(userId, items)
	response, err := c.client.OrderCreate(ctx, request)
	err = detectKnownErrors(err)
	if err != nil {
		return 0, fmt.Errorf("error calling OrderCreate for user %d: %w", userId, err)
	}
	return converter.OrderIdToId(response), nil
}

func (c *Client) GetNumberOfItemInStocks(ctx context.Context, skuId int64) (uint64, error) {
	req := converter.SkuIdToStocksInfoRequest(skuId)
	response, err := c.client.StocksInfo(ctx, req)
	err = detectKnownErrors(err)
	if err != nil {
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