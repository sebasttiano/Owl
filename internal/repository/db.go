package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

// pgError алиас для *pgconn.PgError
var pgError *pgconn.PgError

// ErrNoRows ошибка, если в ответе бд не вернулось ни одной строчки.
var ErrNoRows = errors.New("sql: no rows in result set")

// DBStorage тип реализующий интерфейс service.Repository
type DBStorage struct {
	conn *sqlx.DB
}

// NewDBStorage конструктор для DBStorage
func NewDBStorage(c *sqlx.DB) *DBStorage {
	return &DBStorage{conn: c}
}
