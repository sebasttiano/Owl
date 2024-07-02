package cli

import (
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
	"os"
)

type CLI struct {
	Client *GRPCClient
}

func NewCLI(serverAddr string) *CLI {
	c, err := NewGRPCClient(serverAddr)
	if err != nil {
		logger.Log.Error("failed to init grpc cli", zap.Error(err))
		os.Exit(1)
		return nil
	}
	return &CLI{Client: c}
}
