package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256.ozon.ru/project/loms/internal/controllers/grpc/converter"
	"route256.ozon.ru/project/loms/internal/controllers/grpc/protoc/v1"
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
	v1.UnimplementedLOMServiceServer
	service service
}

func NewServer(service service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) StocksInfo(ctx context.Context, req *v1.StocksInfoRequest) (*v1.StocksInfoResponse, error) {
	skuId := converter.ListStocksInfoRequestToSkuId(req)
	num, err := s.service.GetNumOfAvailable(ctx, skuId)
	if err != nil {
		return nil, err
	}
	return converter.CountToStocksInfoResponse(num), nil
}

func (s *Server) OrderPay(ctx context.Context, orderId *v1.OrderId) (*emptypb.Empty, error) {
	id := converter.OrderIdToId(orderId)
	err := s.service.PayOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *Server) OrderInfo(ctx context.Context, req *v1.OrderId) (*v1.OrderInfoResponse, error) {
	info, err := s.service.GetOrder(ctx, converter.OrderIdToId(req))
	if err != nil {
		return nil, err
	}
	return converter.OrderToOrderInfoResponse(info), nil
}

func (s *Server) OrderCreate(ctx context.Context, req *v1.OrderCreateRequest) (*v1.OrderId, error) {
	userId, items := converter.OrderCreateReqToModel(req)
	orderId, err := s.service.CreateOrder(ctx, userId, items)
	if err != nil {
		return nil, err
	}
	return converter.IdToOrderCreateResponse(orderId), nil
}

func (s *Server) OrderCancel(ctx context.Context, req *v1.OrderId) (*emptypb.Empty, error) {
	err := s.service.CancelOrder(ctx, converter.OrderIdToId(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
