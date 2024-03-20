package controllers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
)

func (s *Server) PayOrder(ctx context.Context, orderId *v1.OrderId) (*emptypb.Empty, error) {
	id := orderId.Id
	err := s.service.PayOrder(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
