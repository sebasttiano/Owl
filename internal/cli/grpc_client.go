package cli

import (
	"errors"
	"fmt"

	"github.com/sebasttiano/Owl/internal/logger"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ErrInitGRPSClient = errors.New("failed to init grpc cli")

// GRPCClient реализующий интерфейс Sender, отправляет на gRPC сервер
type GRPCClient struct {
	Auth   pb.AuthClient
	Binary pb.BinaryClient
	Text   pb.TextClient
	conn   *grpc.ClientConn
}

// NewGRPCClient - конструктор для GRPCClient
func NewGRPCClient(serverAddr string) (*GRPCClient, error) {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("failed to create grpc cli", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrInitGRPSClient, err)
	}
	logger.Log.Info("successfully init grpc cli", zap.String("address", serverAddr))
	auth := pb.NewAuthClient(conn)
	binary := pb.NewBinaryClient(conn)
	text := pb.NewTextClient(conn)

	return &GRPCClient{
		Auth:   auth,
		Binary: binary,
		Text:   text,
		conn:   conn,
	}, nil
}

func (g *GRPCClient) CloseConnection() error {
	err := g.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
