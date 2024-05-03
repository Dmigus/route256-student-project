package http

import (
	"context"
	"errors"
	"net/http"
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
	pprofMux, err := pprofHandler()
	if err != nil {
		return nil, err
	}
	merged.Handle(pprofBasePath, pprofMux)
	gwServer := &http.Server{
		Addr:    addr,
		Handler: merged,
	}
	return &Server{
		serv: gwServer,
	}, nil
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
