package service

import (
	"context"
	"github.com/sebasttiano/Owl/internal/models"
)

type Authenticator interface {
	Register(ctx context.Context, u *models.User) (string, error)
	Login(ctx context.Context, u *models.User) (string, error)
}

func (s *Service) Register(ctx context.Context, u *models.User) (string, error) {
	return "", nil
}

func (s *Service) Login(ctx context.Context, u *models.User) (string, error) {
	return "", nil
}
