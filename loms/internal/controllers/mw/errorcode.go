package mw

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/models"
)

func SetUpErrorCode(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	switch {
	case errors.Is(err, models.ErrNotFound):
		err = status.Error(codes.NotFound, err.Error())
	case errors.Is(err, models.ErrFailedPrecondition):
		err = status.Error(codes.FailedPrecondition, err.Error())
	}
	return resp, err
}
