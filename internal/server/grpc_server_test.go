package server

//import (
//	"context"
//	"github.com/jmoiron/sqlx"
//	"github.com/sebasttiano/Owl/internal/repository"
//	"github.com/sebasttiano/Owl/internal/service"
//	"github.com/stretchr/testify/assert"
//	"net"
//	"sync"
//	"testing"
//	"time"
//)
//
//var GServ *GRPSServer
//
//func init() {
//	conn := &sqlx.DB{}
//	repo, _ := repository.NewDBStorage(conn)
//	srv := service.NewService(&service.Settings{}, repo)
//	GServ = NewGRPSServer(srv)
//}
//
//func TestNewGRPSServer(t *testing.T) {
//	assert.IsType(t, &GRPSServer{}, GServ)
//}
//
//func TestGRPSServer_StartAndShutdown(t *testing.T) {
//
//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
//	defer cancel()
//
//	wg := &sync.WaitGroup{}
//
//	wg.Add(1)
//	go GServ.Start(":4095")
//	time.Sleep(1 * time.Second)
//
//	// Assert socket is used
//	_, err := net.Listen("tcp", ":4095")
//	assert.Error(t, err)
//
//	GServ.HandleShutdown(ctx, wg)
//	wg.Wait()
//
//	// Assert socket is free
//	_, err = net.Listen("tcp", ":4095")
//	assert.NoError(t, err)
//}
