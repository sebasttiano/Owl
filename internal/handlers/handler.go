package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock.go

var (
	ErrInternalGrpc = errors.New("internal grpc server error")
)

var getUserIDFromContext = func(ctx context.Context) (int, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Log.Error("failed get user id from context")
		return 0, status.Errorf(codes.Internal, ErrInternalGrpc.Error())
	}
	userID, err := strconv.Atoi(md.Get("user-id")[0])
	if err != nil {
		logger.Log.Error("failed to convert user id to integer", zap.Error(err))
		return 0, status.Errorf(codes.Internal, ErrInternalGrpc.Error())
	}
	return userID, nil
}

type AuthServer struct {
	Auth     Authenticator
	JManager *JWTManager
	pb.UnimplementedAuthServer
}

type ResourceServer struct {
	Resource ResourceServ
	pb.UnimplementedResourceServer
}

type Authenticator interface {
	Register(ctx context.Context, name, password string) error
	Login(ctx context.Context, name, password string) (int, error)
}

type ResourceServ interface {
	SetResource(ctx context.Context, res models.Resource) (*models.Resource, error)
	GetResource(ctx context.Context, res *models.Resource) (*models.Resource, error)
	GetAllResources(ctx context.Context, uid int) ([]*models.Resource, error)
	DeleteResource(ctx context.Context, res *models.Resource) error
}
