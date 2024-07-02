package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sebasttiano/Owl/internal/logger"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"time"
)

var (
	ErrTokenCreationFailed = errors.New("token creation failed")
	ErrTokenMethod         = errors.New("unexpected token signing method")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidTokenClaims  = errors.New("invalid token claims")
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

// Claims — struct with standart and customs claims
type Claims struct {
	jwt.RegisteredClaims
	ID int `json:"id"`
}

// BuildJWTString creates token and returns it via string.
func (j *JWTManager) BuildJWTString(id int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			Subject:   strconv.Itoa(id),
			Issuer:    "localhost:8080/api/user/login",
			Audience:  []string{"localhost:8080"},
		},
		ID: id,
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenCreationFailed, err)
	}

	return tokenString, nil
}

func (j *JWTManager) Verify(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrTokenMethod
			}

			return []byte(j.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}

func (a *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*emptypb.Empty, error) {
	if err := a.Auth.Register(ctx, in.GetName(), in.GetPassword()); err != nil {
		logger.Log.Error("user register failed", zap.Error(err))
		if errors.Is(err, service.ErrUserAlreadyExist) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, service.ErrUserPasswordSave) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, service.ErrUserRegisrationFailed) {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return nil, nil
}

func (a *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse

	id, err := a.Auth.Login(ctx, in.GetName(), in.GetPassword())
	if err != nil {
		logger.Log.Error("user login failed", zap.Error(err))
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, `user with name %s not found`, in.GetName())
		}
		if errors.Is(err, service.ErrWrongPassword) {
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
		return nil, status.Errorf(codes.Internal, `login failed`)
	}

	token, err := a.JManager.BuildJWTString(id)
	if err != nil {
		logger.Log.Error("token creation failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "token creation failed")
	}

	response.Token = token
	return &response, nil
}
