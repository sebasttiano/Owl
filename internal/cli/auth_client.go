package cli

import (
	"context"
	"time"

	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthClient struct {
	service  pb.AuthClient
	username string
	password string
}

func NewAuthClient(conn *grpc.ClientConn, username string, password string) *AuthClient {
	srv := pb.NewAuthClient(conn)
	return &AuthClient{srv, username, password}
}

func (a *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Name:     a.username,
		Password: a.password,
	}

	res, err := a.service.Login(ctx, req)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return "", service.ErrUserNotFound
			default:
				return "", err
			}
		}
	}

	return res.GetToken(), nil
}
