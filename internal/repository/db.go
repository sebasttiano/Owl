package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Owl/internal/models"
)

// ErrDBNoRows ошибка, если в ответе бд не вернулось ни одной строчки.
var ErrDBNoRows = errors.New("sql: no rows in result set")

// DBStorage тип реализующий интерфейс service.Repository
type DBStorage struct {
	conn *sqlx.DB
}

// NewDBStorage конструктор для DBStorage
func NewDBStorage(c *sqlx.DB) *DBStorage {
	return &DBStorage{conn: c}
}

func (d *DBStorage) GetUserByID(ctx context.Context, user *models.User) error {

	sqlSelect := `SELECT id, name, password FROM users WHERE id = $1`

	if err := d.conn.GetContext(ctx, user, sqlSelect, user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found, %w", ErrDBNoRows)
		}
		return err
	}
	return nil
}

func (d *DBStorage) GetUser(ctx context.Context, user *models.User) error {

	sqlSelect := `SELECT id, name, password FROM users WHERE name = $1`

	if err := d.conn.GetContext(ctx, user, sqlSelect, user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found, %w", ErrDBNoRows)
		}
		return err
	}
	return nil
}

func (d *DBStorage) AddUser(ctx context.Context, user *models.User) error {

	tx, err := d.conn.Beginx()
	if err != nil {
		return err
	}

	// create new user
	sqlInsert := `INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id`
	var id int
	if err := tx.GetContext(ctx, &id, sqlInsert, user.Name, user.HashedPassword); err != nil {
		tx.Rollback()
		return err
	}
	user.ID = id

	tx.Commit()
	return nil
}
