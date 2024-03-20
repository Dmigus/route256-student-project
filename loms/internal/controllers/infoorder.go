package controllers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
	"route256.ozon.ru/project/loms/internal/models"
)

func (s *Server) InfoOrder(ctx context.Context, req *v1.OrderId) (*v1.InfoOrderResponse, error) {
	info, err := s.service.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return orderToInfoOrderResponse(info), nil
}

func orderToInfoOrderResponse(order *models.Order) *v1.InfoOrderResponse {
	response := &v1.InfoOrderResponse{Items: make([]*v1.Item, 0, len(order.Items))}
	response.Status = modelOrderStatusToTransportStatus(order.Status)
	response.User = order.UserId
	for _, it := range order.Items {
		response.Items = append(response.Items, orderItemToItem(it))
	}
	return response
}

func orderItemToItem(item models.OrderItem) *v1.Item {
	return &v1.Item{Sku: uint32(item.SkuId), Count: uint32(item.Count)}
}

func modelOrderStatusToTransportStatus(st models.OrderStatus) v1.OrderStatus {
	switch st {
	case models.New:
		return v1.OrderStatus_ORDER_STATUS_NEW
	case models.AwaitingPayment:
		return v1.OrderStatus_ORDER_STATUS_AWAITING_PAYMENT
	case models.Payed:
		return v1.OrderStatus_ORDER_STATUS_PAYED
	case models.Failed:
		return v1.OrderStatus_ORDER_STATUS_FAILED
	case models.Cancelled:
		return v1.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return v1.OrderStatus_ORDER_STATUS_UNDEFINED
	}
}
