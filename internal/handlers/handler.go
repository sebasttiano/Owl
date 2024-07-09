package handlers

import (
	"context"
	"errors"
	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

var (
	ErrInternalGrpc   = errors.New("internal grpc server error")
	ErrBadRequestGrpc = errors.New("bad request")
)

type AuthServer struct {
	Auth     Authenticator
	JManager *JWTManager
	pb.UnimplementedAuthServer
}

type BinaryServer struct {
	Binary BinaryServ
	pb.UnimplementedBinaryServer
}

type TextServer struct {
	Text TextServ
	pb.UnimplementedTextServer
}

type Authenticator interface {
	Register(ctx context.Context, name, password string) error
	Login(ctx context.Context, name, password string) (int, error)
	Find(ctx context.Context, uid int) (bool, error)
}

type BinaryServ interface {
	SetBinary(ctx context.Context, data models.Resource) error
	GetBinary(ctx context.Context, id int) (models.Resource, error)
	GetAllBinaries(ctx context.Context) ([]models.Resource, error)
	DeleteBinary(ctx context.Context, id int) error
}

type TextServ interface {
	SetText(ctx context.Context, data models.Resource) error
	GetText(ctx context.Context, res *models.Resource) (*models.Resource, error)
	GetAllTexts(ctx context.Context) ([]models.Resource, error)
	DeleteText(ctx context.Context, id int) error
}

type CardServ interface {
	SetCard(ctx context.Context, uid string, data models.Resource) error
	GetCard(ctx context.Context, id string) (models.Resource, error)
	GetAllCards(ctx context.Context, uid string) ([]models.Resource, error)
	DeleteCard(ctx context.Context, id string) error
}

type PasswordServ interface {
	SetPassword(ctx context.Context, uid string, data models.Resource) error
	GetPassword(ctx context.Context, id string) (models.Resource, error)
	GetAllPasswords(ctx context.Context, uid string) ([]models.Resource, error)
	DeletePassword(ctx context.Context, id string) error
}

func getUserIDFromContext(ctx context.Context) (int, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Log.Error("failed get user id from context")
		return 0, status.Errorf(codes.Internal, ErrInternalGrpc.Error())
	}

	userId, err := strconv.Atoi(md.Get("user-id")[0])
	if err != nil {
		logger.Log.Error("failed to convert user id to integer", zap.Error(err))
		return 0, status.Errorf(codes.Internal, ErrInternalGrpc.Error())
	}
	return userId, nil
}
