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
	SetText(ctx context.Context, res *models.ResourceDB, piece *models.PieceDB) error
	GetText(ctx context.Context, res *models.ResourceDB) (*models.ResourceDB, *models.PieceDB, error)
	GetAllTexts(ctx context.Context, uid int) ([]*models.Resource, error)
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

type TextService struct {
	Cipher encrypted.Cipher
	Repo   ResourceRepo
}

func NewTextService(repo ResourceRepo) *TextService {
	return &TextService{Repo: repo, Cipher: encrypted.CFBCipher{}}
}
