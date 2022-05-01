package repository

import (
	"context"
	"errors"
	"rodnik/internal/apperror"
	"sync"
	"time"
)

type tokenMemory struct {
	sync.Mutex
	db map[string]string
}

func NewTokenMemory() *tokenMemory {
	return &tokenMemory{
		db: make(map[string]string),
	}
}

func (r *tokenMemory) CreateToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	r.Lock()
	defer r.Unlock()
	if r.db == nil {
		r.db = make(map[string]string)
	}
	r.db[tokenID] = userID
	return nil
}
func (r *tokenMemory) GetUserIDByToken(ctx context.Context, tokenID string) (string, error) {
	userId, ok := r.db[tokenID]
	if !ok {
		return "", apperror.Authorization.New(ErrorMessageTokenNotFound)
	}
	return userId, nil
}
func (r *tokenMemory) DeleteToken(ctx context.Context, tokenID string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.db[tokenID]; !ok {
		//todo Переформулировать сообщение об ошибке. Вынести ошибки в отдельный файл
		return errors.New("Could not found refresh token")
	}
	delete(r.db, tokenID)
	return nil
}
func (r *tokenMemory) DeleteUserTokens(ctx context.Context, userId string) error {
	r.Lock()
	defer r.Unlock()
	for key, value := range r.db {
		if value != userId {
			continue
		}
		delete(r.db, key)
	}
	return nil
}
