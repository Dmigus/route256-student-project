package loms

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256.ozon.ru/project/cart/internal/clients/loms/converter"
	v1 "route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"
	"route256.ozon.ru/project/cart/internal/models"
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
	if err != nil {
		return 0, err
	}
	return converter.OrderIdToId(response), nil
}

func (c *Client) GetNumberOfItemInStocks(ctx context.Context, skuId int64) (uint64, error) {
	req := converter.SkuIdToStocksInfoRequest(skuId)
	response, err := c.client.StocksInfo(ctx, req)
	if err != nil {
		return 0, err
	}
	return converter.StocksInfoResponseToCount(response), nil
}
