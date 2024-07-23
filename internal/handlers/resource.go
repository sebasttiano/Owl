package handlers

import (
	"context"
	"errors"

	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *ResourceServer) SetResource(ctx context.Context, in *pb.SetResourceRequest) (*pb.SetResourceResponse, error) {

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resource := models.Resource{
		UserID:      userID,
		Type:        models.ResourceType(in.Resource.Type),
		Description: in.Resource.Description,
		Content:     in.Resource.Content,
	}

	resp, err := t.Resource.SetResource(ctx, resource)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal grpc server error")
	}
	response := &pb.SetResourceResponse{Resource: &pb.ResourceMeta{Id: int32(resp.ID), Description: resp.Description}}
	return response, nil
}

func (t *ResourceServer) GetResource(ctx context.Context, in *pb.GetResourceRequest) (*pb.GetResourceResponse, error) {

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res := &models.Resource{ID: int(in.Id), UserID: userID}
	resource, err := t.Resource.GetResource(ctx, res)
	if err != nil {
		if errors.Is(err, service.ErrGetResourceNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal grpc server error")
	}

	req := pb.GetResourceResponse{Resource: &pb.ResourceMsg{Content: resource.Content, Description: resource.Description, Type: string(resource.Type)}}
	return &req, nil

}

func (t *ResourceServer) GetAllResources(ctx context.Context, _ *emptypb.Empty) (*pb.GetAllResourcesResponse, error) {

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	resources, err := t.Resource.GetAllResources(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal grpc server error")
	}
	resourcesMeta := make([]*pb.ResourceMeta, len(resources))
	for i, resource := range resources {
		resourcesMeta[i] = &pb.ResourceMeta{Id: int32(resource.ID), Description: resource.Description, Type: string(resource.Type)}
	}
	return &pb.GetAllResourcesResponse{Resources: resourcesMeta}, nil
}

func (t *ResourceServer) DeleteResource(ctx context.Context, in *pb.DeleteResourceRequest) (*emptypb.Empty, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	res := &models.Resource{ID: int(in.Id), UserID: userID}
	if err := t.Resource.DeleteResource(ctx, res); err != nil {
		return nil, status.Error(codes.Internal, "internal grpc server error")
	}

	return nil, nil
}
