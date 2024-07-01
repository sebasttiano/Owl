package handlers

import (
	"context"
	"github.com/sebasttiano/Owl/internal/models"
)

type Authenticator interface {
	Register(ctx context.Context, u *models.User) (string, error)
	Login(ctx context.Context, u *models.User) (string, error)
}

type BinaryServ interface {
	SetBinary(ctx context.Context, uid string, data models.Resource) error
	GetBinary(ctx context.Context, id string) (models.Resource, error)
	GetAllBinaries(ctx context.Context, uid string) ([]models.Resource, error)
	DeleteBinary(ctx context.Context, id string) error
}

type TextServ interface {
	SetText(ctx context.Context, uid string, data models.Resource) error
	GetText(ctx context.Context, id string) (models.Resource, error)
	GetAllTexts(ctx context.Context, uid string) ([]models.Resource, error)
	DeleteText(ctx context.Context, id string) error
}

type CardServ interface {
	SetCard(ctx context.Context, uid string, data models.Resource) error
	GetCard(ctx context.Context, id string) (models.Resource, error)
	GetAllCards(ctx context.Context, uid string) ([]models.Resource, error)
	DeleteCard(ctx context.Context, id string) error
}

type PasswordServ interface {
	SetPassword(ctx context.Context, uid string, data models.Resource) error
	GetPassword(ctx context.Context, id string) (models.Resource, error)
	GetAllPasswords(ctx context.Context, uid string) ([]models.Resource, error)
	DeletePassword(ctx context.Context, id string) error
}
