package repository

import (
	"context"
	"github.com/google/uuid"
	"rodnik/internal/entity"
	"time"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user *entity.User) error
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	//TODO Добавить полное обновление пользователя по всем полям?
	UpdateUserBalance(ctx context.Context, user *entity.User) error
	FindById(ctx context.Context, userId uuid.UUID) (*entity.User, error)
}

type RefreshTokens interface {
	CreateToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	GetUserIDByToken(ctx context.Context, tokenID string) (string, error)
	DeleteToken(ctx context.Context, tokenID string) error
	DeleteUserTokens(ctx context.Context, userId string) error
}

type Tasks interface {
	Create(ctx context.Context, task *entity.Task) error
}
