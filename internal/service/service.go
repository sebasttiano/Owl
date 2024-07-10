package service

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sebasttiano/Owl/internal/encrypted"
	"github.com/sebasttiano/Owl/internal/models"
)

// pgError алиас для *pgconn.PgError
var pgError *pgconn.PgError

type AuthRepo interface {
	GetUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, user *models.User) error
	AddUser(ctx context.Context, user *models.User) error
}

type ResourceRepo interface {
	GetUserHashPass(ctx context.Context, uid int) (string, error)
	SetResource(ctx context.Context, res *models.Resource, piece *models.Piece) (*models.Resource, error)
	GetResource(ctx context.Context, res *models.Resource) (*models.Resource, *models.Piece, error)
	DelResource(ctx context.Context, res *models.Resource) error
	GetAllResources(ctx context.Context, uid int) ([]*models.Resource, error)
}

type AuthService struct {
	Repo AuthRepo
}

func NewAuthService(repo AuthRepo) *AuthService {
	return &AuthService{Repo: repo}
}

type BinaryService struct {
	Repo ResourceRepo
}

func NewBinaryService(repo ResourceRepo) *BinaryService {
	return &BinaryService{Repo: repo}
}

type ResourceService struct {
	Cipher encrypted.Cipher
	Repo   ResourceRepo
}

func NewTextService(repo ResourceRepo) *ResourceService {
	return &ResourceService{Repo: repo, Cipher: encrypted.CFBCipher{}}
}
