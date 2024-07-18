package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sebasttiano/Owl/internal/encrypted"
	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/models"
	"go.uber.org/zap"
)

var (
	ErrSetResource         = errors.New("failed to save resource")
	ErrGetResource         = errors.New("failed to get resource")
	ErrGetResourceNotFound = errors.New("not found requested resource")
	ErrGetAllResources     = errors.New("failed to get all resources meta")
	ErrDelResource         = errors.New("failed to delete resource")
)

func (t *ResourceService) SetResource(ctx context.Context, data models.Resource) (*models.Resource, error) {

	m := &models.Meta{Type: data.Type, Description: data.Description}
	meta, err := json.Marshal(m)
	if err != nil {
		logger.Log.Error("failed to marshal to json", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrSetResource, err)
	}
	res := &models.Resource{UserID: data.UserID, Meta: string(meta), Type: data.Type}

	userPass, err := t.Repo.GetUserHashPass(ctx, data.UserID)
	if err != nil {
		logger.Log.Error("failed get hashed error pass", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrSetResource, err)
	}

	piece, err := encrypted.PasswordEncryption(userPass, t.Cipher, data.Content)
	if err != nil {
		logger.Log.Error("failed to encrypt incoming data", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrSetResource, err)
	}
	resp, err := t.Repo.SetResource(ctx, res, piece)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *ResourceService) GetResource(ctx context.Context, res *models.Resource) (*models.Resource, error) {

	res, piece, err := t.Repo.GetResource(ctx, res)
	if err != nil {
		logger.Log.Error("failed to get resource from repo", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrGetResourceNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrGetResource, err)
	}

	userPass, err := t.Repo.GetUserHashPass(ctx, res.UserID)
	if err != nil {
		logger.Log.Error("failed get hashed error pass", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetResource, err)
	}

	content, err := encrypted.PasswordDecryption(userPass, t.Cipher, piece)
	if err != nil {
		logger.Log.Error("failed to decrypted data", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetResource, err)
	}

	var m models.Meta
	if err = json.Unmarshal([]byte(res.Meta), &m); err != nil {
		logger.Log.Error("failed to unmarshall meta", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetResource, err)
	}

	data := &models.Resource{ID: res.ID, UserID: res.UserID, Type: res.Type, Content: content, Description: m.Description}

	return data, nil
}

func (t *ResourceService) GetAllResources(ctx context.Context, uid int) ([]*models.Resource, error) {

	resources, err := t.Repo.GetAllResources(ctx, uid)
	if err != nil {
		logger.Log.Error("failed to get resources from repo", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetAllResources, err)
	}
	for _, res := range resources {
		res.SetDescriptionFromMeta()
	}
	return resources, nil
}

func (t *ResourceService) DeleteResource(ctx context.Context, res *models.Resource) error {
	if err := t.Repo.DelResource(ctx, res); err != nil {
		logger.Log.Error("failed to delete resource from repo", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrDelResource, err)
	}
	return nil
}
