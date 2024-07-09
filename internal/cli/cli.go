package cli

import (
	"context"
	"errors"
	"github.com/sebasttiano/Owl/internal/config"
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	ErrInitAuthConn        = errors.New("failed to create auth connection")
	ErrInitAuthInterceptor = errors.New("failed to create auth interceptor")
	ErrInitCLIConn         = errors.New("failed to create cli connection")
	ErrInitGRPSClient      = errors.New("failed to create grpc client")
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

	name, pass, err := c.GetUserCreds(ctx)

	authConn, err := grpc.NewClient(c.cfg.GetServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
		textMethodsPath     = "/main.Text/"
		binaryMethodsPath   = "/main.Binary/"
		cardMethodsPath     = "/main.Card/"
		passwordMethodsPath = "/main.Password/"
	)

	return map[string]bool{
		textMethodsPath + "SetText":             true,
		textMethodsPath + "GetText":             true,
		textMethodsPath + "GetAllTexts":         true,
		textMethodsPath + "DeleteText":          true,
		binaryMethodsPath + "SetBinary":         true,
		binaryMethodsPath + "GetBinary":         true,
		binaryMethodsPath + "GetAllBinaries":    true,
		binaryMethodsPath + "DeleteBinary":      true,
		cardMethodsPath + "SetCard":             true,
		cardMethodsPath + "GetCard":             true,
		cardMethodsPath + "GetAllCards":         true,
		cardMethodsPath + "DeleteCard":          true,
		passwordMethodsPath + "SetPassword":     true,
		passwordMethodsPath + "GetPassword":     true,
		passwordMethodsPath + "GetAllPasswords": true,
		passwordMethodsPath + "DeletePassword":  true,
	}
}
