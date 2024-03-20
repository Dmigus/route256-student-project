package controllers

import (
	"context"
	servicepb "route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
	"route256.ozon.ru/project/loms/internal/models"
)

type service interface {
	CreateOrder(ctx context.Context, userId int64, items []models.OrderItem) (int64, error)
	PayOrder(ctx context.Context, orderId int64) error
	GetNumOfAvailable(ctx context.Context, skuId int64) (uint64, error)
	GetOrder(ctx context.Context, orderId int64) (*models.Order, error)
	CancelOrder(ctx context.Context, orderId int64) error
}

type Server struct {
	servicepb.UnimplementedLOMServiceServer
	service service
}

func NewServer(service service) *Server {
	return &Server{
		service: service,
	}
}
