package service

import (
	"context"
	"github.com/sebasttiano/Owl/internal/models"
)

func (b *BinaryService) SetBinary(ctx context.Context, data models.Resource) error {
	return nil
}

func (b *BinaryService) GetBinary(ctx context.Context, id int) (models.Resource, error) {
	return models.Resource{}, nil
}

func (b *BinaryService) GetAllBinaries(ctx context.Context) ([]models.Resource, error) {
	return []models.Resource{}, nil
}

func (b *BinaryService) DeleteBinary(ctx context.Context, id int) error {
	return nil
}
