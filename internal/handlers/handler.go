package handlers

import (
	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
)

type KeeperServer struct {
	Service service.KeeperService
	pb.UnimplementedKeeperServer
}
