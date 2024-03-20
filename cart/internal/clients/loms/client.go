package loms

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"
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
