package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

var ErrNoMetadata = errors.New("metadata is not provided")
var ErrNoAccessToken = errors.New("authorization token is not provided")

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			logger.Info(msg)
		}
	})
}

type AuthInterceptor struct {
	jwtManager       *JWTManager
	whitelistMethods map[string]bool
}

func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, map[string]bool{
		"/main.Auth/Login":    true,
		"/main.Auth/Register": true,
	}}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Log.Info(fmt.Sprintf("--> unary interceptor: %s", info.FullMethod))
		newCtx, err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		if newCtx != nil {
			ctx = newCtx
		}
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	fmt.Println(method)
	_, ok := i.whitelistMethods[method]
	if ok {
		// everyone can access
		return nil, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, ErrNoMetadata.Error())
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, ErrNoAccessToken.Error())
	}

	accessToken := values[0]
	claims, err := i.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%w: %v", ErrInvalidToken, err)
	}

	// Set user id to metadata
	md, ok = metadata.FromIncomingContext(ctx)
	if ok {
		md.Append("user-id", strconv.Itoa(claims.ID))
	}

	newCtx := metadata.NewIncomingContext(ctx, md)
	return newCtx, nil
}
