package repository

import (
	"context"
	"rodnik/internal/entity"
	"time"
)

type Users interface {
	Create(ctx context.Context, user *entity.User) error
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
}

type RefreshTokens interface {
	CreateToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	GetUserIDByToken(ctx context.Context, tokenID string) (string, error)
	DeleteToken(ctx context.Context, tokenID string) error
	DeleteUserTokens(ctx context.Context, userId string) error
}
