package application

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Owl/internal/config"
	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/repository"
	"github.com/sebasttiano/Owl/internal/server"
	"github.com/sebasttiano/Owl/internal/service"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Run() {

	cfg, err := config.NewServerConfig()
	if err != nil {
		fmt.Println("parsing config failed")
		return
	}

	if err := logger.Initialize(cfg.Logger.Level); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	var conn *sqlx.DB
	conn, err = sqlx.Connect("pgx", cfg.GetDSN())
	if err != nil {
		logger.Log.Error("database openning failed", zap.Error(err))
		os.Exit(1)
	}
	defer conn.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	wg := &sync.WaitGroup{}

	repo := repository.NewDBStorage(conn)
	serviceSettings := &service.ServiceSettings{}
	srv := service.NewService(repo, serviceSettings)
	grpcSrv := server.NewGRPSServer(srv)

	wg.Add(1)
	go grpcSrv.Start(cfg.GetAddress())
	go grpcSrv.HandleShutdown(ctx, wg)
	wg.Wait()

}
