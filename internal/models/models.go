package models

const (
	ResourceTypeBinary   = "BINARY"
	ResourceTypeText     = "TEXT"
	ResourceTypeCard     = "CARD"
	ResourceTypePassword = "PASSWORD"
)

type User struct {
	ID             int
	Name           string `json:"name" valid:"required,type(string)"`
	HashedPassword string `json:"password" db:"password" valid:"required,type(string)"`
	RegisteredAT   string
}

type Resource struct {
	ID      string
	PieceID string
	UserID  string
	Type    string
	Meta    string
}
