package controllers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
)

func (s *Server) CancelOrder(ctx context.Context, req *v1.OrderId) (*emptypb.Empty, error) {
	err := s.service.CancelOrder(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
