package main

import (
	"fmt"
	"github.com/sebasttiano/Owl/internal/cli"
	"github.com/sebasttiano/Owl/internal/config"
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize("DEBUG"); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	cfg, err := config.NewClientConfig()
	if err != nil {
		logger.Log.Error("parsing config failed", zap.Error(err))
		return
	}

	cliApp := cli.NewCLI(cfg)
	cliApp.Run()

	//if err != nil {
	//	if e, ok := status.FromError(err); ok {
	//		print(e.Code(), "\n")
	//		print(e.Message(), "\n")
	//	}
	//	os.Exit(1)
	//}
}
