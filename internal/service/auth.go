package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/sebasttiano/Owl/internal/logger"
	"github.com/sebasttiano/Owl/internal/models"
	"github.com/sebasttiano/Owl/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserRegisrationFailed = errors.New("user registration failed")
	ErrUserPasswordSave      = errors.New("failed to save password. check your input")
	ErrUserLoginFailed       = errors.New("user login failed")
	ErrWrongPassword         = errors.New("wrong password")
	ErrUserAlreadyExist      = errors.New("user already exist")
)

func (a *AuthService) Register(ctx context.Context, name, password string) error {
	hashedPass, err := a.hashPassword(password)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUserPasswordSave, err)
	}
	user := models.User{Name: name, HashedPassword: hashedPass}
	if err := a.Repo.AddUser(ctx, &user); err != nil {
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return ErrUserAlreadyExist
			}
		}
		return ErrUserRegisrationFailed
	}
	logger.Log.Info("registered new user: " + name)
	return nil
}

func (a *AuthService) Login(ctx context.Context, name, password string) (int, error) {

	user := models.User{Name: name}
	if err := a.Repo.GetUser(ctx, &user); err != nil {
		if errors.Is(err, repository.ErrDBNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("%w: %v", ErrUserLoginFailed, err)
	}

	// Check password
	if err := a.checkPasswordHash(password, user.HashedPassword); err != nil {
		return 0, ErrWrongPassword
	}
	return user.ID, nil
}

func (a *AuthService) Find(ctx context.Context, uid int) (bool, error) {
	var user models.User
	user.ID = uid
	err := a.Repo.GetUserByID(ctx, &user)
	if err != nil {
		if errors.Is(err, repository.ErrDBNoRows) {
			return false, ErrUserNotFound
		}
		return false, err
	}
	return true, nil
}

func (a *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (a *AuthService) checkPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrWrongPassword
	}
	return nil
}
