package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Owl/internal/repository"

	"github.com/sebasttiano/Owl/internal/handlers"
	"github.com/sebasttiano/Owl/internal/logger"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPSServerSettings struct {
	SecretKey     string
	CertFile      string
	CertKey       string
	TokenDuration time.Duration
}

// GRPSServer реалиузет gRPC сервер.
type GRPSServer struct {
	srv *grpc.Server
}

// NewGRPSServer конструктор для gRPC сервера
func NewGRPSServer(conn *sqlx.DB, settings *GRPSServerSettings) *GRPSServer {
	repo := repository.NewDBStorage(conn)
	j := handlers.NewJWTManager(settings.SecretKey, settings.TokenDuration)
	authInterceptor := handlers.NewAuthInterceptor(j)

	tlsCredentials, err := loadTLSCredentials(settings.CertFile, settings.CertKey)
	if err != nil {
		fmt.Println(err)
		logger.Log.Fatal("cannot load TLS credentials", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(handlers.InterceptorLogger(logger.Log)),
			authInterceptor.Unary(),
		),
	)
	pb.RegisterAuthServer(s, &handlers.AuthServer{
		Auth:     service.NewAuthService(repo),
		JManager: j,
	})
	pb.RegisterResourceServer(s, &handlers.ResourceServer{Resource: service.NewTextService(repo)})
	return &GRPSServer{
		srv: s,
	}
}

// Start запускает grpc сервер.
func (s *GRPSServer) Start(addr string) {
	logger.Log.Info("Running gRPC server", zap.String("address", addr))
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Log.Error("failed to allocate tcp socket for gRPC server", zap.Error(err))
		os.Exit(1)
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
