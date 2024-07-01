package service

import (
	"context"
	"github.com/sebasttiano/Owl/internal/models"
)

func (t *TextService) SetText(ctx context.Context, uid string, data models.Resource) error {
	return nil
}

func (t *TextService) GetText(ctx context.Context, id string) (models.Resource, error) {
	return models.Resource{}, nil
}

func (t *TextService) GetAllTexts(ctx context.Context, uid string) ([]models.Resource, error) {
	return []models.Resource{}, nil
}

func (t *TextService) DeleteText(ctx context.Context, id string) error {
	return nil
}
