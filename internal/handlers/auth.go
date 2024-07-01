package handlers

import pb "github.com/sebasttiano/Owl/internal/proto"

type KeeperServer struct {
	Auth   Authenticator
	Binary BinaryServ
	Text   TextServ
	pb.UnimplementedAuthServer
	pb.UnimplementedBinaryServer
	pb.UnimplementedTextServer
}
