package handlers

import (
	"context"
	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *TextServer) SetText(ctx context.Context, in *pb.SetTextRequest) (*pb.SetTextResponse, error) {

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

	resp, err := t.Text.SetText(ctx, resource)
	if err != nil {
		return nil, err
	}
	response := &pb.SetTextResponse{Text: &pb.TextMeta{Id: int32(resp.ID), Description: resp.Description}}
	return response, nil
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

func (t *TextServer) GetAllTexts(ctx context.Context, _ *emptypb.Empty) (*pb.GetAllTextsResponse, error) {

	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	texts, err := t.Text.GetAllTexts(ctx, userId)
	if err != nil {
		return nil, err
	}

	textsMeta := make([]*pb.TextMeta, len(texts))
	for i, text := range texts {
		textsMeta[i] = &pb.TextMeta{Id: int32(text.ID), Description: text.Description}
	}
	return &pb.GetAllTextsResponse{Texts: textsMeta}, nil
}

func (t *TextServer) DeleteText(ctx context.Context, in *pb.DeleteTextRequest) (*emptypb.Empty, error) {
	return nil, nil
}
