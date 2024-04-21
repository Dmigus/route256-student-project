package http

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	errorsPkg "github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
	"time"
)

const (
	basePath    = "/api/"
	metricsPath = "/metrics"
)

// Server это gateway к grpc контроллеру loms
type Server struct {
	serv *http.Server
}

// NewServer инициализирует Server с добавленным swagger и метриками
func NewServer(lomsaddress, addr, swaggerPath string, metricsHandler http.Handler) (*Server, error) {
	gwmux, err := initGateWayMux(lomsaddress)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

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

// Serve это блокирующий вызов, который запускает http контроллер обработатывать входящих запросов
func (s *Server) Serve() error {
	if err := s.serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop останавливает запущенный http контроллер в течение timeout времени
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
