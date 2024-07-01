package service

import (
	"context"
	"github.com/sebasttiano/Owl/internal/models"
)

func (b *BinaryService) SetBinary(ctx context.Context, uid string, data models.Resource) error {
	return nil
}

func (b *BinaryService) GetBinary(ctx context.Context, id string) (models.Resource, error) {
	return models.Resource{}, nil
}

func (b *BinaryService) GetAllBinaries(ctx context.Context, uid string) ([]models.Resource, error) {
	return []models.Resource{}, nil
}

func (b *BinaryService) DeleteBinary(ctx context.Context, id string) error {
	return nil
}
