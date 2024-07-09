package handlers

import (
	"context"
	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *TextServer) SetText(ctx context.Context, in *pb.SetTextRequest) (*emptypb.Empty, error) {

	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resource := models.Resource{
		UserID:      userId,
		Type:        models.Text,
		Description: in.Text.Description,
		Content:     in.Text.Text,
	}
	t.Text.SetText(ctx, resource)

	return nil, nil
}

func (t *TextServer) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {

	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res := &models.Resource{ID: int(in.Id), UserID: userId}
	text, err := t.Text.GetText(ctx, res)
	if err != nil {
		return nil, err
	}

	req := pb.GetTextResponse{Text: &pb.TextMsg{Text: text.Content, Description: text.Description}}
	return &req, nil

}

func (t *TextServer) GetAllTexts(ctx context.Context, in *emptypb.Empty) (*pb.GetAllTextsResponse, error) {
	return nil, nil
}

func (t *TextServer) DeleteText(ctx context.Context, in *pb.DeleteTextRequest) (*emptypb.Empty, error) {
	return nil, nil
}
