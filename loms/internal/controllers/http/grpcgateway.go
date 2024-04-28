package http

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	errorsPkg "github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
)

func initGateWayMux(lomsaddress string) (*runtime.ServeMux, error) {
	conn, err := grpc.Dial(lomsaddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errorsPkg.Wrap(err, "failed to dial")
	}
	gwmux := runtime.NewServeMux(runtime.WithErrorHandler(fixFailedPreconditionCodeMapping))
	if err = v1.RegisterLOMServiceHandler(context.Background(), gwmux, conn); err != nil {
		return nil, errorsPkg.Wrap(err, "failed to register gateway")
	}
	return gwmux, nil
}

func fixFailedPreconditionCodeMapping(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	gRPCCode := status.Code(err)
	if gRPCCode == codes.FailedPrecondition {
		err = &runtime.HTTPStatusError{
			HTTPStatus: http.StatusPreconditionFailed,
			Err:        err,
		}
	}
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, err)
}
