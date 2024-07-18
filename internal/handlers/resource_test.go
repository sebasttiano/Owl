package handlers

import (
	"context"
	"github.com/golang/mock/gomock"
	mockservice "github.com/sebasttiano/Owl/internal/handlers/mocks"
	"github.com/sebasttiano/Owl/internal/models"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func TestResourceServer_SetResource(t *testing.T) {
	type mockBehaviour func(r *mockservice.MockResourceServ, in *pb.SetResourceRequest)

	testTable := []struct {
		name          string
		in            *pb.SetResourceRequest
		mockBehaviour mockBehaviour
		expected      *pb.SetResourceResponse
		err           error
	}{
		{
			name: "OK set text",
			in: &pb.SetResourceRequest{Resource: &pb.ResourceMsg{
				Type:        string(models.Text),
				Content:     "Zhili bili ded and babka",
				Description: "Kolobok",
			}},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.SetResourceRequest) {
				response := &models.Resource{UserID: 1, Type: models.ResourceType(in.Resource.Type), Content: in.Resource.Content,
					ID: 10, Description: in.Resource.Description}
				r.EXPECT().SetResource(gomock.Any(), models.Resource{UserID: 1, Type: models.ResourceType(in.Resource.Type), Content: in.Resource.Content, Description: in.Resource.Description}).Return(response, nil)
			},
			expected: &pb.SetResourceResponse{
				Resource: &pb.ResourceMeta{Id: 10, Description: "Kolobok"},
			},
			err: nil,
		},
		{
			name: "OK set bank",
			in: &pb.SetResourceRequest{Resource: &pb.ResourceMsg{
				Type:        string(models.Card),
				Content:     "'ccn': '3343 4890 5543 1290', 'exp': '04/25', 'cvv': '109', 'holder': 'Ivanov Ivan'",
				Description: "MY BANK",
			}},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.SetResourceRequest) {
				response := &models.Resource{UserID: 1, Type: models.ResourceType(in.Resource.Type), Content: in.Resource.Content,
					ID: 10, Description: in.Resource.Description}
				r.EXPECT().SetResource(gomock.Any(), models.Resource{UserID: 1, Type: models.ResourceType(in.Resource.Type), Content: in.Resource.Content, Description: in.Resource.Description}).Return(response, nil)
			},
			expected: &pb.SetResourceResponse{
				Resource: &pb.ResourceMeta{Id: 10, Description: "MY BANK"},
			},
			err: nil,
		},
		{
			name: "OK set creds",
			in: &pb.SetResourceRequest{Resource: &pb.ResourceMsg{
				Type:        string(models.Password),
				Content:     "'username': 'ivan', 'password': 'supersecret'",
				Description: "mail",
			}},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.SetResourceRequest) {
				response := &models.Resource{UserID: 1, Type: models.ResourceType(in.Resource.Type), Content: in.Resource.Content,
					ID: 10, Description: in.Resource.Description}
				r.EXPECT().SetResource(gomock.Any(), models.Resource{UserID: 1, Type: models.ResourceType(in.Resource.Type), Content: in.Resource.Content, Description: in.Resource.Description}).Return(response, nil)
			},
			expected: &pb.SetResourceResponse{
				Resource: &pb.ResourceMeta{Id: 10, Description: "mail"},
			},
			err: nil,
		},
		{
			name: "NOT OK set resource",
			in: &pb.SetResourceRequest{Resource: &pb.ResourceMsg{
				Type:        string(models.Password),
				Content:     "'username': 'ivan', 'password': 'supersecret'",
				Description: "my first pass",
			}},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.SetResourceRequest) {
				r.EXPECT().SetResource(gomock.Any(), models.Resource{
					UserID: 1, Type: models.ResourceType(in.Resource.Type), Content: in.Resource.Content, Description: in.Resource.Description,
				}).Return(nil, service.ErrSetResource)
			},
			expected: nil,
			err:      status.Errorf(codes.Internal, "internal grpc server error"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			getUserIDFromContext = func(ctx context.Context) (int, error) {
				return 1, nil
			}

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()
			mock := mockservice.NewMockResourceServ(c)
			tt.mockBehaviour(mock, tt.in)
			pb.RegisterResourceServer(s, &ResourceServer{Resource: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()
			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}
			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewResourceClient(conn)

			resp, err := client.SetResource(ctx, tt.in)
			if tt.err != nil {
				require.Errorf(t, err, tt.err.Error())
				require.Equal(t, tt.err, err)
			} else {
				assert.Equal(t, tt.expected.Resource.Id, resp.Resource.Id)
				assert.Equal(t, tt.expected.Resource.Description, resp.Resource.Description)
			}
		})
	}
}

func TestResourceServer_GetResource(t *testing.T) {
	type mockBehaviour func(r *mockservice.MockResourceServ, in *pb.GetResourceRequest)

	testTable := []struct {
		name          string
		in            *pb.GetResourceRequest
		mockBehaviour mockBehaviour
		expected      *pb.GetResourceResponse
		err           error
	}{
		{
			name: "OK get text",
			in:   &pb.GetResourceRequest{Id: 10},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.GetResourceRequest) {
				res := models.Resource{ID: int(in.Id)}
				respRes := res
				res.UserID = 1
				respRes.Description = "Kolobok"
				respRes.Content = "Zhili bili ded and babka"
				respRes.Type = models.Text
				r.EXPECT().GetResource(gomock.Any(), &res).Return(&respRes, nil)
			},
			expected: &pb.GetResourceResponse{Resource: &pb.ResourceMsg{Type: string(models.Text), Content: "Zhili bili ded and babka",
				Description: "Kolobok"}},
			err: nil,
		},
		{
			name: "OK get bank",
			in:   &pb.GetResourceRequest{Id: 9},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.GetResourceRequest) {
				res := models.Resource{ID: int(in.Id)}
				respRes := res
				res.UserID = 1
				respRes.Description = "MY BANK"
				respRes.Content = "'ccn': '3343 4890 5543 1290', 'exp': '04/25', 'cvv': '109', 'holder': 'Ivanov Ivan'"
				respRes.Type = models.Card
				r.EXPECT().GetResource(gomock.Any(), &res).Return(&respRes, nil)
			},
			expected: &pb.GetResourceResponse{Resource: &pb.ResourceMsg{Type: string(models.Card), Content: "'ccn': '3343 4890 5543 1290', 'exp': '04/25', 'cvv': '109', 'holder': 'Ivanov Ivan'",
				Description: "MY BANK"}},
			err: nil,
		},
		{
			name: "OK get creds",
			in:   &pb.GetResourceRequest{Id: 9},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.GetResourceRequest) {
				res := models.Resource{ID: int(in.Id)}
				respRes := res
				res.UserID = 1
				respRes.Description = "my first pass"
				respRes.Content = "'username': 'ivan', 'password': 'supersecret'"
				respRes.Type = models.Password
				r.EXPECT().GetResource(gomock.Any(), &res).Return(&respRes, nil)
			},
			expected: &pb.GetResourceResponse{Resource: &pb.ResourceMsg{Type: string(models.Password), Content: "'username': 'ivan', 'password': 'supersecret'",
				Description: "my first pass"}},
			err: nil,
		},
		{
			name: "NOT OK set resource. NOT FOUND",
			in:   &pb.GetResourceRequest{Id: 9},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.GetResourceRequest) {
				res := models.Resource{ID: int(in.Id)}
				res.UserID = 1
				r.EXPECT().GetResource(gomock.Any(), &res).Return(nil, service.ErrGetResourceNotFound)
			},
			expected: nil,
			err:      status.Errorf(codes.NotFound, "not found requested resource"),
		},
		{
			name: "NOT OK set resource. INTERNAL ERROR",
			in:   &pb.GetResourceRequest{Id: 9},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.GetResourceRequest) {
				res := models.Resource{ID: int(in.Id)}
				res.UserID = 1
				r.EXPECT().GetResource(gomock.Any(), &res).Return(nil, service.ErrGetResource)
			},
			expected: nil,
			err:      status.Errorf(codes.Internal, "internal grpc server error"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			getUserIDFromContext = func(ctx context.Context) (int, error) {
				return 1, nil
			}

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()
			mock := mockservice.NewMockResourceServ(c)
			tt.mockBehaviour(mock, tt.in)
			pb.RegisterResourceServer(s, &ResourceServer{Resource: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()
			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}
			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewResourceClient(conn)
			resp, err := client.GetResource(ctx, tt.in)
			if tt.err != nil {
				require.Errorf(t, err, tt.err.Error())
				require.Equal(t, tt.err, err)
			} else {
				assert.Equal(t, tt.expected.Resource.Content, resp.Resource.Content)
				assert.Equal(t, tt.expected.Resource.Description, resp.Resource.Description)
				assert.Equal(t, tt.expected.Resource.Type, resp.Resource.Type)
			}
		})
	}
}

func TestResourceServer_GetAllResources(t *testing.T) {
	type mockBehaviour func(r *mockservice.MockResourceServ, in *emptypb.Empty)

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		expected      *pb.GetAllResourcesResponse
		err           error
	}{
		{
			name: "OK get all resources",
			mockBehaviour: func(r *mockservice.MockResourceServ, in *emptypb.Empty) {
				res := []*models.Resource{
					{
						Type:        models.Text,
						ID:          9,
						Description: "Kolobok",
						Content:     "Zhili bili ded and babka",
					},
					{
						Type:        models.Card,
						ID:          10,
						Description: "MY BANK",
						Content:     "'ccn': '3343 4890 5543 1290', 'exp': '04/25', 'cvv': '109', 'holder': 'Ivanov Ivan'",
					},
					{
						Type:        models.Password,
						ID:          11,
						Description: "my first pass",
						Content:     "'username': 'ivan', 'password': 'supersecret'",
					},
				}
				r.EXPECT().GetAllResources(gomock.Any(), 1).Return(res, nil)
			},
			expected: &pb.GetAllResourcesResponse{Resources: []*pb.ResourceMeta{
				{Id: 9, Description: "Kolobok", Type: string(models.Text)},
				{Id: 10, Description: "MY BANK", Type: string(models.Card)},
				{Id: 11, Description: "my first pass", Type: string(models.Password)},
			}},
			err: nil,
		},
		{
			name: "NOT OK get all resources",
			mockBehaviour: func(r *mockservice.MockResourceServ, in *emptypb.Empty) {
				r.EXPECT().GetAllResources(gomock.Any(), 1).Return(nil, service.ErrGetAllResources)
			},
			expected: nil,
			err:      status.Errorf(codes.Internal, "internal grpc server error"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			getUserIDFromContext = func(ctx context.Context) (int, error) {
				return 1, nil
			}

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()
			mock := mockservice.NewMockResourceServ(c)
			tt.mockBehaviour(mock, &emptypb.Empty{})
			pb.RegisterResourceServer(s, &ResourceServer{Resource: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()
			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}
			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewResourceClient(conn)
			resp, err := client.GetAllResources(ctx, &emptypb.Empty{})
			if tt.err != nil {
				require.Errorf(t, err, tt.err.Error())
				require.Equal(t, tt.err, err)
			} else {
				for i, res := range tt.expected.Resources {
					assert.Equal(t, res.Id, resp.Resources[i].Id)
					assert.Equal(t, res.Description, resp.Resources[i].Description)
					assert.Equal(t, res.Type, resp.Resources[i].Type)
				}
			}
		})
	}
}

func TestResourceServer_DeleteResource(t *testing.T) {
	type mockBehaviour func(r *mockservice.MockResourceServ, in *pb.DeleteResourceRequest)

	testTable := []struct {
		name          string
		in            *pb.DeleteResourceRequest
		mockBehaviour mockBehaviour
		err           error
	}{
		{
			name: "OK delete resource",
			in:   &pb.DeleteResourceRequest{Id: 10},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.DeleteResourceRequest) {
				res := &models.Resource{ID: int(in.Id), UserID: 1}
				r.EXPECT().DeleteResource(gomock.Any(), res).Return(nil)
			},
			err: nil,
		},
		{
			name: "NOT OK delete resource",
			in:   &pb.DeleteResourceRequest{Id: 10},
			mockBehaviour: func(r *mockservice.MockResourceServ, in *pb.DeleteResourceRequest) {
				res := &models.Resource{ID: int(in.Id), UserID: 1}
				r.EXPECT().DeleteResource(gomock.Any(), res).Return(service.ErrDelResource)
			},
			err: status.Errorf(codes.Internal, "internal grpc server error"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			getUserIDFromContext = func(ctx context.Context) (int, error) {
				return 1, nil
			}

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()
			mock := mockservice.NewMockResourceServ(c)
			tt.mockBehaviour(mock, tt.in)
			pb.RegisterResourceServer(s, &ResourceServer{Resource: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()
			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}
			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewResourceClient(conn)
			_, err = client.DeleteResource(ctx, tt.in)
			if tt.err != nil {
				require.Errorf(t, err, tt.err.Error())
				require.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
