package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	converter2 "route256.ozon.ru/project/loms/internal/controllers/grpc/converter"
	v12 "route256.ozon.ru/project/loms/internal/controllers/grpc/protoc/v1"
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
	v12.UnimplementedLOMServiceServer
	service service
}

func NewServer(service service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) StocksInfo(ctx context.Context, req *v12.StocksInfoRequest) (*v12.StocksInfoResponse, error) {
	skuId := converter2.ListStocksInfoRequestToSkuId(req)
	num, err := s.service.GetNumOfAvailable(ctx, skuId)
	if err != nil {
		return nil, err
	}
	return converter2.CountToStocksInfoResponse(num), nil
}

func (s *Server) OrderPay(ctx context.Context, orderId *v12.OrderId) (*emptypb.Empty, error) {
	id := converter2.OrderIdToId(orderId)
	err := s.service.PayOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *Server) OrderInfo(ctx context.Context, req *v12.OrderId) (*v12.OrderInfoResponse, error) {
	info, err := s.service.GetOrder(ctx, converter2.OrderIdToId(req))
	if err != nil {
		return nil, err
	}
	return converter2.OrderToOrderInfoResponse(info), nil
}

func (s *Server) OrderCreate(ctx context.Context, req *v12.OrderCreateRequest) (*v12.OrderId, error) {
	userId, items := converter2.OrderCreateReqToModel(req)
	orderId, err := s.service.CreateOrder(ctx, userId, items)
	if err != nil {
		return nil, err
	}
	return converter2.IdToOrderCreateResponse(orderId), nil
}

func (s *Server) OrderCancel(ctx context.Context, req *v12.OrderId) (*emptypb.Empty, error) {
	err := s.service.CancelOrder(ctx, converter2.OrderIdToId(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
