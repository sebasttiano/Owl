package server

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"testing"
	"time"
)

var GServ *GRPSServer

func init() {
	conn := &sqlx.DB{}
	//repo := repository.NewDBStorage(conn)
	//srv := service.NewTextService(repo)
	GServ = NewGRPSServer(conn, &GRPSServerSettings{})
}

func TestNewGRPSServer(t *testing.T) {
	assert.IsType(t, &GRPSServer{}, GServ)
}

func TestGRPSServer_StartAndShutdown(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go GServ.Start(":4095")
	time.Sleep(1 * time.Second)

	// Assert socket is used
	_, err := net.Listen("tcp", ":4095")
	require.Error(t, err)

	GServ.HandleShutdown(ctx, wg)
	wg.Wait()

	// Assert socket is free
	_, err = net.Listen("tcp", ":4095")
	assert.NoError(t, err)
}
