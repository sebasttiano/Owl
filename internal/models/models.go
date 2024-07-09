package models

import (
	"encoding/json"
	"fmt"
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
)

type ResourceType string

const (
	Binary   ResourceType = "BINARY"
	Text     ResourceType = "TEXT"
	Card     ResourceType = "CARD"
	Password ResourceType = "PASSWORD"
)

type (
	User struct {
		ID             int
		Name           string `json:"name" valid:"required,type(string)"`
		HashedPassword string `json:"password" db:"password" valid:"required,type(string)"`
		RegisteredAT   string
	}

	Resource struct {
		ID          int          `db:"id"`
		UserID      int          `db:"user_id"`
		PieceUUID   string       `db:"piece_uuid"`
		Meta        string       `db:"meta"`
		Type        ResourceType `db:"type"`
		Description string
		Content     string
	}
	// Piece is a piece of encrypted information.
	Piece struct {
		Content []byte `db:"content"` // Content of the piece.
		IV      []byte `db:"iv"`
		Salt    []byte `db:"salt"`
		Meta    string // Meta info of the piece.
	}

	Meta struct {
		Type        ResourceType `json:"type"`
		Description string       `json:"description"`
	}
)

func (r *Resource) SetDescriptionFromMeta() {

	var m Meta
	if err := json.Unmarshal([]byte(r.Meta), &m); err != nil {
		logger.Log.Error(fmt.Sprintf("failed to unmarshall meta: %s", r.Meta), zap.Error(err))
		return
	}
	r.Description = m.Description
}
