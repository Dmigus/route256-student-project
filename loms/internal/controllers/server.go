package controllers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Server) CreateOrder(ctx context.Context, req *servicepb.CreateOrderRequest) (*servicepb.CreateOrderResponse, error) {
	userId, items := reqToModel(req)
	orderId, err := s.service.CreateOrder(ctx, userId, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return idToCreateOrderResponse(orderId), nil
}

func reqToModel(req *servicepb.CreateOrderRequest) (int64, []models.OrderItem) {
	items := make([]models.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, itemToOrderItem(it))
	}
	return req.User, items
}

func itemToOrderItem(it *servicepb.Item) models.OrderItem {
	return models.OrderItem{SkuId: int64(it.Sku), Count: uint16(it.Count)}
}

func idToCreateOrderResponse(id int64) *servicepb.CreateOrderResponse {
	return &servicepb.CreateOrderResponse{OrderID: id}
}

//InfoOrder(context.Context, *v1.InfoOrderRequest) (*v1.InfoOrderResponse, error)
//PayOrder(context.Context, *v1.PayOrderRequest) (*emptypb.Empty, error)
//CancelOrder(context.Context, *v1.CancelOrderRequest) (*emptypb.Empty, error)
//ListStocksInfo(context.Context, *v1.ListStocksInfoRequest) (*v1.ListStocksInfoResponse, error)
