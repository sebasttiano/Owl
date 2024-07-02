package service

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sebasttiano/Owl/internal/models"
)

// pgError алиас для *pgconn.PgError
var pgError *pgconn.PgError

type Repository interface {
	GetUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, user *models.User) error
	AddUser(ctx context.Context, user *models.User) error
}

type AuthService struct {
	Repo Repository
}

func NewAuthService(repo Repository) *AuthService {
	return &AuthService{Repo: repo}
}

type BinaryService struct {
	Repo *Repository
}

func NewBinaryService(repo *Repository) *BinaryService {
	return &BinaryService{Repo: repo}
}

type TextService struct {
	Repo *Repository
}

func NewTextService(repo *Repository) *TextService {
	return &TextService{Repo: repo}
}
