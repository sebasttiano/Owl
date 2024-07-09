package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sebasttiano/Owl/internal/encrypted"
	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/models"
	"go.uber.org/zap"
)

var (
	ErrSetText = errors.New("failed to save text")
	ErrGetText = errors.New("failed to get text")
)

func (t *TextService) SetText(ctx context.Context, data models.Resource) error {

	m := &models.Meta{Type: data.Type, Description: data.Description}
	meta, err := json.Marshal(m)
	//meta, err := json.Marshal(
	//	map[string]any{
	//		"type":        data.Type,
	//		"description": data.Description,
	//	},
	//)
	if err != nil {
		logger.Log.Error("failed to marshal to json", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrSetText, err)
	}
	res := &models.ResourceDB{UserID: data.UserID, Meta: string(meta), Type: data.Type}
	//content := make([]byte, len(data.Content))
	//fmt.Println(data.Content)
	//fmt.Println([]byte(data.Content))
	//if _, err := base64.RawStdEncoding.Decode(content, []byte(data.Content)); err != nil {
	//	logger.Log.Error("failed to decode to base64", zap.Error(err))
	//	return fmt.Errorf("%w: %v", ErrSetText, err)
	//}

	userPass, err := t.Repo.GetUserHashPass(ctx, data.UserID)
	if err != nil {
		logger.Log.Error("failed get hashed error pass", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrSetText, err)
	}

	piece, err := encrypted.PasswordEncryption(userPass, t.Cipher, data.Content)
	if err != nil {
		logger.Log.Error("failed to encrypt incoming data", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrSetText, err)
	}

	//reader := cipher.StreamReader{
	//	S: t.Cipher.Encrypter(enc.Block, enc.IV),
	//	R: bytes.NewReader([]byte(data.Content)),
	//}
	//
	//content, err := io.ReadAll(reader)
	//if err != nil {
	//	logger.Log.Error("failed to cipher content", zap.Error(err))
	//	return fmt.Errorf("%w: %v", ErrSetText, err)
	//}

	//piece := &models.PieceDB{Content: encryptedContent, IV: enc.IV, Salt: enc.Salt}
	if err := t.Repo.SetText(ctx, res, piece); err != nil {
		return err
	}
	return nil
}

func (t *TextService) GetText(ctx context.Context, resource *models.Resource) (*models.Resource, error) {

	res := &models.ResourceDB{ID: resource.ID}
	res, piece, err := t.Repo.GetText(ctx, res)
	if err != nil {
		logger.Log.Error("failed to get resource from repo", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetText, err)
	}

	userPass, err := t.Repo.GetUserHashPass(ctx, resource.UserID)
	if err != nil {
		logger.Log.Error("failed get hashed error pass", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetText, err)
	}

	content, err := encrypted.PasswordDecryption(userPass, t.Cipher, piece)
	if err != nil {
		logger.Log.Error("failed to decrypted data", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetText, err)
	}

	var m models.Meta
	if err = json.Unmarshal([]byte(res.Meta), &m); err != nil {
		logger.Log.Error("failed to unmarshall meta", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrGetText, err)
	}

	data := &models.Resource{ID: res.ID, UserID: res.UserID, Type: res.Type, Content: content, Description: m.Description}

	return data, nil
}

func (t *TextService) GetAllTexts(ctx context.Context) ([]models.Resource, error) {
	return []models.Resource{}, nil
}

func (t *TextService) DeleteText(ctx context.Context, id int) error {
	return nil
}
