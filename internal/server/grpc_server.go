package server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sebasttiano/Owl/internal/handlers"
	"github.com/sebasttiano/Owl/internal/logger"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GRPSServer реалиузет gRPC сервер.
type GRPSServer struct {
	srv *grpc.Server
}

// NewGRPSServer конструктор для gRPC сервера
func NewGRPSServer(repo service.Repository) *GRPSServer {
	s := grpc.NewServer(grpc.UnaryInterceptor(logging.UnaryServerInterceptor(handlers.InterceptorLogger(logger.Log))))
	pb.RegisterAuthServer(s, &handlers.KeeperServer{
		Auth:   service.NewAuthService(&repo),
		Binary: service.NewBinaryService(&repo),
		Text:   service.NewTextService(&repo),
	})
	return &GRPSServer{
		srv: s,
	}
}

// Start запускает grpc сервер.
func (s *GRPSServer) Start(addr string) {
	logger.Log.Info("Running gRPC server", zap.String("address", addr))
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		logger.Log.Error("failed to allocate tcp socket for gRPC server", zap.Error(err))
	}
	if err := s.srv.Serve(listen); err != nil {
		logger.Log.Error("failed to start gRPC server", zap.Error(err))
	}
}

// HandleShutdown закрывает grpc сервер.
func (s *GRPSServer) HandleShutdown(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	<-ctx.Done()
	logger.Log.Info("shutdown signal caught. shutting down gRPC server")

	s.srv.GracefulStop()
	logger.Log.Info("gRPC server gracefully shutdown")
}
