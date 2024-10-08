// Package encryption provides encryption utils.
package encrypted

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/models"
	"go.uber.org/zap"

	"golang.org/x/crypto/pbkdf2"
)

// Data is encryption data.
type Data struct {
	Block cipher.Block
	IV    []byte
	Salt  []byte
}

// PasswordEncryption returns password based encryption data.
func PasswordEncryption(password string, ciph Cipher, content []byte) (*models.Piece, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	block, blockError := aes.NewCipher(
		pbkdf2.Key(([]byte)(password), salt, 4096, 32, sha256.New),
	)
	if blockError != nil {
		return nil, blockError
	}

	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	reader := cipher.StreamReader{
		S: ciph.Encrypter(block, iv),
		R: bytes.NewReader(content),
	}

	encryptedContent, err := io.ReadAll(reader)
	if err != nil {
		logger.Log.Error("failed to cipher content", zap.Error(err))
		return nil, err
	}

	piece := &models.Piece{Content: encryptedContent, IV: iv, Salt: salt}
	return piece, nil
}

// PasswordDecryption return decrypted password
func PasswordDecryption(password string, ciph Cipher, piece *models.Piece) ([]byte, error) {
	block, err := aes.NewCipher(
		pbkdf2.Key([]byte(password), piece.Salt, keyIter, keyLen, sha256.New),
	)
	if err != nil {
		return nil, err
	}

	reader := cipher.StreamReader{
		S: ciph.Decrypter(block, piece.IV),
		R: bytes.NewReader(piece.Content),
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return content, nil
}
