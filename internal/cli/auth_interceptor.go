package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sebasttiano/Owl/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var ErrRefreshToken = errors.New("failed to get refresh")

type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(
	authClient *AuthClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (i *AuthInterceptor) refreshToken() error {
	accessToken, err := i.authClient.Login()
	if err != nil {
		return err
	}

	i.accessToken = accessToken
	logger.Log.Info(fmt.Sprintf("token refreshed: %v", accessToken))
	return nil
}

func (i *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := i.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := i.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}

func (i *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		//logger.Log.Debug(fmt.Sprintf("--> unary interceptor: %s", method))
		if i.authMethods[method] {
			return invoker(i.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", i.accessToken)
}
