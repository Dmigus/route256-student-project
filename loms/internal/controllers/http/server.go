package http

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
	"time"
)

const (
	basePath    = "/api/"
	metricsPath = "/metrics"
)

type Server struct {
	serv *http.Server
}

func NewServer(lomsaddress, addr, swaggerPath string, metricsHandler http.Handler) *Server {
	gwmux := initGateWayMux(lomsaddress)
	swaggeruimux := swaggerUIHandler(swaggerPath)
	merged := http.NewServeMux()
	merged.Handle(basePath, gwmux)
	merged.Handle(swaggeruiprefix, swaggeruimux)
	if metricsHandler != nil {
		merged.Handle(metricsPath, metricsHandler)
	}
	gwServer := &http.Server{
		Addr:    addr,
		Handler: merged,
	}
	return &Server{
		serv: gwServer,
	}
}

func initGateWayMux(lomsaddress string) *runtime.ServeMux {
	conn, err := grpc.Dial(lomsaddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial:", err)
	}
	gwmux := runtime.NewServeMux(runtime.WithErrorHandler(fixFailedPreconditionCodeMapping))
	if err = v1.RegisterLOMServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	return gwmux
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

func (s *Server) Serve() error {
	if err := s.serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(timeout time.Duration) error {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()
	if err := s.serv.Shutdown(shutdownCtx); err != nil {
		errClose := s.serv.Close()
		err = errors.Join(err, errClose)
		return err
	}
	return nil
}
