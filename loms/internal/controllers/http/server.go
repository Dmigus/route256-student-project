package http

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	v1 "route256.ozon.ru/project/loms/internal/controllers/grpc/protoc/v1"
	"time"
)

type Server struct {
	serv *http.Server
}

func NewServer(lomsaddress, addr string) *Server {
	conn, err := grpc.Dial(lomsaddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial:", err)
	}
	gwmux := runtime.NewServeMux()
	if err = v1.RegisterLOMServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    addr,
		Handler: gwmux,
	}
	return &Server{
		serv: gwServer,
	}
}

func (s *Server) Serve() {
	if err := s.serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (s *Server) Stop(timeout time.Duration) {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()
	if err := s.serv.Shutdown(shutdownCtx); err != nil {
		err = errors.Join(err, s.serv.Close())
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}
