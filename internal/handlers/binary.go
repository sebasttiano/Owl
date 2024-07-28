package handlers

import (
	"context"
	"fmt"
	pb "github.com/sebasttiano/Owl/internal/proto"
)

func (b *BinaryServer) SetBinary(ctx context.Context, in *pb.SetBinaryRequest) (*pb.SetBinaryResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(userID)
	return nil, nil
}
