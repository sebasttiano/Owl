package models

type ResourceType string

const (
	Binary   ResourceType = "BINARY"
	Text     ResourceType = "TEXT"
	Card     ResourceType = "CARD"
	Password ResourceType = "PASSWORD"
)

type (
	// Piece is a piece of encrypted information.
	Piece struct {
		Content []byte // Content of the piece.
		Meta    string // Meta info of the piece.
	}
	User struct {
		ID             int
		Name           string `json:"name" valid:"required,type(string)"`
		HashedPassword string `json:"password" db:"password" valid:"required,type(string)"`
		RegisteredAT   string
	}

	Resource struct {
		ID          int
		UserID      int
		Type        ResourceType
		Description string
		Content     string
	}

	ResourceDB struct {
		ID        int          `db:"id"`
		UserID    int          `db:"user_id"`
		PieceUUID string       `db:"piece_uuid"`
		Meta      string       `db:"meta"`
		Type      ResourceType `db:"type"`
	}
	PieceDB struct {
		Content []byte `db:"content"`
		IV      []byte `db:"iv"`
		Salt    []byte `db:"salt"`
	}
	BlobDB struct {
		location string
	}

	Meta struct {
		Type        ResourceType `json:"type"`
		Description string       `json:"description"`
	}
)
