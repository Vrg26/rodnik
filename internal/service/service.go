package service

import (
	"context"
	"rodnik/internal/entity"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Token interface {
	GetTokenPair(ctx context.Context, userId string) (*TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	ParseToken(tokenString string) (*CustomClaims, error)
	DeleteUserTokens(ctx context.Context, userId string) error
}

type Users interface {
	Create(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, user *entity.User) error
}

type Tasks interface {
	Create(ctx context.Context, task *entity.Task) error
}
