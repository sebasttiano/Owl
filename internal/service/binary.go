package service

import (
	"context"

	"github.com/sebasttiano/Owl/internal/models"
)

func (b *BinaryService) SetBinary(_ context.Context, _ models.Resource) error {
	return nil
}

func (b *BinaryService) GetBinary(_ context.Context, _ int) (models.Resource, error) {
	return models.Resource{}, nil
}

func (b *BinaryService) GetAllBinaries(_ context.Context) ([]models.Resource, error) {
	return []models.Resource{}, nil
}

func (b *BinaryService) DeleteBinary(_ context.Context, _ int) error {
	return nil
}
