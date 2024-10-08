package cli

import (
	"context"
	"errors"
	"time"

	"github.com/sebasttiano/Owl/internal/config"
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	ErrInitAuthConn        = errors.New("failed to create auth connection")
	ErrInitAuthInterceptor = errors.New("failed to create auth interceptor")
	ErrInitCLIConn         = errors.New("failed to create cli connection")
	ErrInitGRPSClient      = errors.New("failed to create grpc client")
)

var (
	msgHeight,
	msgWidth int
)

type CLI struct {
	Auth   *AuthClient
	Client *GRPCClient
	cfg    *config.ClientConfig
}

func NewCLI(cfg *config.ClientConfig) *CLI {
	return &CLI{cfg: cfg}
}

func (c *CLI) Run() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tlsCredentials, err := loadTLSCredentials(c.cfg.Cert.CA)
	if err != nil {
		logger.Log.Fatal("cannot load TLS credentials: ", zap.Error(err))
	}

	name, pass, err := c.GetUserCreds(ctx, tlsCredentials)
	if err != nil {
		if errors.Is(err, ErrCanceledByUser) {
			logger.Log.Info(ErrCanceledByUser.Error())
			return nil
		}
		logger.Log.Error("failed to get user creds", zap.Error(err))
		return err
	}

	authConn, err := grpc.NewClient(c.cfg.GetServerAddress(), grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		logger.Log.Error("failed to create auth connection", zap.Error(err))
		return ErrInitAuthConn
	}

	authClient := NewAuthClient(authConn, name, pass)
	logger.Log.Info("successfully init auth client", zap.String("address", authConn.Target()))

	authInterceptor, err := NewAuthInterceptor(authClient, AuthMethods(), time.Duration(c.cfg.Auth.RefreshPeriod)*time.Second)
	if err != nil {
		logger.Log.Error("failed to create auth interceptor", zap.Error(err))
		return ErrInitAuthInterceptor
	}

	conn, err := grpc.NewClient(
		c.cfg.GetServerAddress(),
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()))

	if err != nil {
		logger.Log.Error("failed to create cli connection", zap.Error(err))
		return ErrInitCLIConn
	}

	c.Client, err = NewGRPCClient(conn)
	if err != nil {
		logger.Log.Error("failed to create grpc client")
		return ErrInitGRPSClient
	}

	if err := c.StartMainBoard(ctx); err != nil {
		return err
	}

	return nil
}

func AuthMethods() map[string]bool {
	const (
		resourceMethodsPath = "/main.Resource/"
		binaryMethodsPath   = "/main.Binary/"
	)

	return map[string]bool{
		resourceMethodsPath + "SetResource":     true,
		resourceMethodsPath + "GetResource":     true,
		resourceMethodsPath + "GetAllResources": true,
		resourceMethodsPath + "DeleteResource":  true,
		binaryMethodsPath + "SetBinary":         true,
		binaryMethodsPath + "GetBinary":         true,
		binaryMethodsPath + "GetAllBinaries":    true,
		binaryMethodsPath + "DeleteBinary":      true,
	}
}
