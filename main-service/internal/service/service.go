package service

import (
	"context"
	"main-service/internal/entity"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Token interface {
	GetTokenPair(ctx context.Context, userId string) (*TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	ParseToken(tokenString string) (*CustomClaims, error)
	DeleteUserTokens(ctx context.Context, userId string) error
}

type Users interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Login(ctx context.Context, user *entity.User) error
	SetAvatar(ctx context.Context, userID string, imageBytes []byte) (string, error)
}

type Tasks interface {
	Create(ctx context.Context, task *entity.Task) (*entity.Task, error)
}
