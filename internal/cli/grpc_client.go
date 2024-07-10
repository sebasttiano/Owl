package cli

import (
	"context"
	"time"

	pb "github.com/sebasttiano/Owl/internal/proto"
	"google.golang.org/grpc"
)

// GRPCClient реализующий интерфейс Sender, отправляет на gRPC сервер
type GRPCClient struct {
	Auth     pb.AuthClient
	Binary   pb.BinaryClient
	Resource pb.ResourceClient
	conn     *grpc.ClientConn
	username string
	password string
}

// NewGRPCClient - конструктор для GRPCClient
func NewGRPCClient(conn *grpc.ClientConn) (*GRPCClient, error) {
	// устанавливаем соединение с сервером
	auth := pb.NewAuthClient(conn)
	binary := pb.NewBinaryClient(conn)
	text := pb.NewResourceClient(conn)

	return &GRPCClient{
		Auth:     auth,
		Binary:   binary,
		Resource: text,
		conn:     conn,
	}, nil
}

func (g *GRPCClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Name:     g.username,
		Password: g.password,
	}

	res, err := g.Auth.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetToken(), nil
}

func (g *GRPCClient) CloseConnection() error {
	err := g.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
