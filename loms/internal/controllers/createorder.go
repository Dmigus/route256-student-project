package controllers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
	"route256.ozon.ru/project/loms/internal/models"
)

func (s *Server) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.OrderId, error) {
	userId, items := createOrderReqToModel(req)
	orderId, err := s.service.CreateOrder(ctx, userId, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return idToCreateOrderResponse(orderId), nil
}

func createOrderReqToModel(req *v1.CreateOrderRequest) (int64, []models.OrderItem) {
	items := make([]models.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, itemToOrderItem(it))
	}
	return req.User, items
}

func itemToOrderItem(it *v1.Item) models.OrderItem {
	return models.OrderItem{SkuId: int64(it.Sku), Count: uint16(it.Count)}
}

func idToCreateOrderResponse(id int64) *v1.OrderId {
	return &v1.OrderId{Id: id}
}
