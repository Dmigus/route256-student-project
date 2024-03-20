package controllers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
)

func (s *Server) ListStocksInfo(ctx context.Context, req *v1.ListStocksInfoRequest) (*v1.ListStocksInfoResponse, error) {
	skuId := listStocksInfoRequestToSkuId(req)
	num, err := s.service.GetNumOfAvailable(ctx, skuId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &v1.ListStocksInfoResponse{Count: num}, nil
}

func listStocksInfoRequestToSkuId(req *v1.ListStocksInfoRequest) int64 {
	return int64(req.Sku)
}
