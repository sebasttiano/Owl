package main

import (
	"context"
	"fmt"
	"github.com/sebasttiano/Owl/internal/cli"
	"github.com/sebasttiano/Owl/internal/logger"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"time"
)

func main() {
	if err := logger.Initialize("DEBUG"); err != nil {
		fmt.Println("logger initialization failed")
		return
	}
	cli := cli.NewCLI(":8091")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cli.Client.Text.SetText(ctx, &pb.SetTextRequest{Text: &pb.TextMsg{Text: "Some text", Description: "test"}})
	fmt.Println(cli)
}
