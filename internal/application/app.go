package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Owl/internal/config"
	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/server"
	"go.uber.org/zap"
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

	settings := &server.GRPSServerSettings{
		SecretKey:     cfg.Server.Secret,
		TokenDuration: time.Duration(cfg.Server.TokenDuration) * time.Second,
		CertFile:      cfg.Cert.Cert,
		CertKey:       cfg.Cert.Key,
	}

	grpcSrv := server.NewGRPSServer(conn, settings)

	wg.Add(1)
	go grpcSrv.Start(cfg.GetAddress())
	go grpcSrv.HandleShutdown(ctx, wg)
	wg.Wait()

}
