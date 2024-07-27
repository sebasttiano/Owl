package handlers

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	mockservice "github.com/sebasttiano/Owl/internal/handlers/mocks"
	pb "github.com/sebasttiano/Owl/internal/proto"
	"github.com/sebasttiano/Owl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestAuthServer_Login(t *testing.T) {
	type mockBehavior func(r *mockservice.MockAuthenticator, in *pb.LoginRequest)

	testTable := []struct {
		name         string
		in           *pb.LoginRequest
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "OK Login",
			in:   &pb.LoginRequest{Name: "Father", Password: "My pass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.LoginRequest) {
				r.EXPECT().Login(gomock.Any(), in.GetName(), in.GetPassword()).Return(1, nil)
			},
			err: nil,
		},
		{
			name: "NOT OK Login. not found",
			in:   &pb.LoginRequest{Name: "Father", Password: "my pass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.LoginRequest) {
				r.EXPECT().Login(gomock.Any(), in.GetName(), in.GetPassword()).Return(0, service.ErrUserNotFound)
			},
			err: status.Error(codes.NotFound, "user with name Father not found"),
		},
		{
			name: "NOT OK Login. wrong password",
			in:   &pb.LoginRequest{Name: "Father", Password: "wrong pass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.LoginRequest) {
				r.EXPECT().Login(gomock.Any(), in.GetName(), in.GetPassword()).Return(0, service.ErrWrongPassword)
			},
			err: status.Error(codes.Aborted, service.ErrWrongPassword.Error()),
		},
		{
			name: "NOT OK Login. internal",
			in:   &pb.LoginRequest{Name: "Father", Password: "my pass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.LoginRequest) {
				r.EXPECT().Login(gomock.Any(), in.GetName(), in.GetPassword()).Return(0, errors.New("unexpected error"))
			},
			err: status.Error(codes.Internal, "login failed"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()
			mock := mockservice.NewMockAuthenticator(c)
			tt.mockBehavior(mock, tt.in)
			pb.RegisterAuthServer(s, &AuthServer{Auth: mock, JManager: NewJWTManager("secret", 3)})
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
			client := pb.NewAuthClient(conn)
			resp, err := client.Login(ctx, tt.in)
			if tt.err != nil {
				require.Errorf(t, err, tt.err.Error())
				require.Equal(t, tt.err, err)
			} else {
				assert.Equal(t, 212, len(resp.Token))
			}
		})
	}

}

func TestAuthServer_Register(t *testing.T) {
	type mockBehavior func(r *mockservice.MockAuthenticator, in *pb.RegisterRequest)

	testTable := []struct {
		name         string
		in           *pb.RegisterRequest
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "OK register",
			in:   &pb.RegisterRequest{Name: "Father", Password: "my pass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.RegisterRequest) {
				r.EXPECT().Register(gomock.Any(), in.GetName(), in.GetPassword()).Return(nil)
			},
			err: nil,
		},
		{
			name: "NOT OK register. already exist",
			in:   &pb.RegisterRequest{Name: "Father", Password: "my pass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.RegisterRequest) {
				r.EXPECT().Register(gomock.Any(), in.GetName(), in.GetPassword()).Return(service.ErrUserAlreadyExist)
			},
			err: status.Error(codes.AlreadyExists, service.ErrUserAlreadyExist.Error()),
		},
		{
			name: "NOT OK register. invalid argument",
			in:   &pb.RegisterRequest{Name: "Father", Password: "myвыфвфыpass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.RegisterRequest) {
				r.EXPECT().Register(gomock.Any(), in.GetName(), in.GetPassword()).Return(service.ErrUserPasswordSave)
			},
			err: status.Error(codes.InvalidArgument, service.ErrUserPasswordSave.Error()),
		},
		{
			name: "NOT OK register. internal error",
			in:   &pb.RegisterRequest{Name: "Father", Password: "my pass"},
			mockBehavior: func(r *mockservice.MockAuthenticator, in *pb.RegisterRequest) {
				r.EXPECT().Register(gomock.Any(), in.GetName(), in.GetPassword()).Return(service.ErrUserRegisrationFailed)
			},
			err: status.Error(codes.Internal, service.ErrUserRegisrationFailed.Error()),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()
			mock := mockservice.NewMockAuthenticator(c)
			tt.mockBehavior(mock, tt.in)
			pb.RegisterAuthServer(s, &AuthServer{Auth: mock, JManager: NewJWTManager("secret", 3)})
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
			client := pb.NewAuthClient(conn)
			_, err = client.Register(ctx, tt.in)
			if tt.err != nil {
				require.Errorf(t, err, tt.err.Error())
				require.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
