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

func (d *DBStorage) GetUserHashPass(ctx context.Context, uid int) (string, error) {

	//type
	var pass *string
	sqlSelect := `SELECT password FROM users WHERE id = $1`

	if err := d.conn.GetContext(ctx, &pass, sqlSelect, uid); err != nil {
		return "", err
	}
	return *pass, nil
}

func (d *DBStorage) SetResource(ctx context.Context, res *models.Resource, piece *models.Piece) (*models.Resource, error) {

	tx, err := d.conn.Beginx()
	if err != nil {
		return nil, err
	}

	// insert into resources
	sqlInsert := `INSERT INTO resources (user_id, type, meta) VALUES ($1, $2, $3) RETURNING id`
	var id int
	if err := d.conn.GetContext(ctx, &id, sqlInsert, res.UserID, res.Type, res.Meta); err != nil {
		tx.Rollback()
		return nil, err
	}

	// insert into pieces
	sqlInsert = `INSERT INTO pieces (content, iv, salt, resource_id) VALUES ($1, $2, $3, $4)`
	if _, err := d.conn.ExecContext(ctx, sqlInsert, piece.Content, piece.IV, piece.Salt, id); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	res.ID = id
	return res, nil
}

func (d *DBStorage) GetResource(ctx context.Context, res *models.Resource) (*models.Resource, *models.Piece, error) {

	// get resource
	sqlSelect := `SELECT type, meta FROM resources WHERE id = $1 AND user_id = $2`
	if err := d.conn.GetContext(ctx, res, sqlSelect, res.ID, res.UserID); err != nil {
		return nil, nil, err
	}

	// get piece
	piece := &models.Piece{}
	sqlSelect = `SELECT content, salt, iv FROM pieces WHERE resource_id = $1`
	if err := d.conn.GetContext(ctx, piece, sqlSelect, res.ID); err != nil {
		return nil, nil, err
	}
	return res, piece, nil
}
func (d *DBStorage) GetAllResources(ctx context.Context, uid int) ([]*models.Resource, error) {

	var resources []*models.Resource
	// get resources
	sqlSelect := `SELECT id, meta, type FROM resources WHERE user_id = $1`
	if err := d.conn.SelectContext(ctx, &resources, sqlSelect, uid); err != nil {
		return nil, err
	}
	return resources, nil
}

func (d *DBStorage) DelResource(ctx context.Context, res *models.Resource) error {

	tx, err := d.conn.Beginx()
	if err != nil {
		return err
	}
	// del resource
	sqlDelete := `DELETE FROM resources WHERE id = $1 AND user_id = $2`
	_, err = d.conn.ExecContext(ctx, sqlDelete, res.ID, res.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
