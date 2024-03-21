package mw

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/providers/orders"
	"route256.ozon.ru/project/loms/internal/providers/stocks"
	"route256.ozon.ru/project/loms/internal/usecases"
)

func InterpretErrorToCode(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else if errors.Is(err, stocks.ErrItemIsNotExists) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else if errors.Is(err, stocks.ErrNotEnoughItems) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		} else if errors.Is(err, models.ErrWrongOrderStatus) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		} else if errors.Is(err, usecases.ErrService) {
			return nil, status.Errorf(codes.Internal, usecases.ErrService.Error())
		} else {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return resp, nil
}
