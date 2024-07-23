package server

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var GServ *GRPSServer

func init() {
	conn := &sqlx.DB{}
	GServ = NewGRPSServer(conn, &GRPSServerSettings{CertKey: "../../cert/server-key.pem", CertFile: "../../cert/server-cert.pem"})
}

func TestNewGRPSServer(t *testing.T) {
	assert.IsType(t, &GRPSServer{}, GServ)
}

func TestGRPSServer_StartAndShutdown(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
