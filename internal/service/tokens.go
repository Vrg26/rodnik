package service

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"rodnik/internal/apperror"
	"rodnik/internal/repository"
	"rodnik/pkg/logger"
	"time"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// todo вынести в слой моделей
type TokenPair struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken uuid.UUID `json:"refreshToken"`
}

type TokenService struct {
	Repo                  repository.RefreshTokens
	l                     logger.Logger
	SecretKey             []byte
	ExpirationTimeSec     int64
	RefreshExpirationTime int64
}

func NewTokenService(repo repository.RefreshTokens, l logger.Logger, sk []byte, et int64, ret int64) *TokenService {
	return &TokenService{repo, l, sk, et, ret}
}

func (s *TokenService) GetTokenPair(ctx context.Context, userId string) (*TokenPair, error) {
	unixTimeNow := time.Now().Unix()
	tokenExp := unixTimeNow + s.ExpirationTimeSec

	claims := &CustomClaims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExp,
			IssuedAt:  unixTimeNow,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(s.SecretKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	err = s.Repo.CreateToken(ctx, userId, refreshToken.String(), time.Duration(s.RefreshExpirationTime)*time.Minute)
	if err != nil {
		return nil, err
	}

	return &TokenPair{accessToken, refreshToken}, nil
}

func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	userID, err := s.Repo.GetUserIDByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if err = s.Repo.DeleteToken(ctx, refreshToken); err != nil {
		return nil, err
	}
	return s.GetTokenPair(ctx, userID)
}

func (s *TokenService) ParseToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, apperror.Authorization.New(ErrorMessageInvalidToken)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, apperror.Authorization.New(ErrorMessageFailedConvertToClaims)
	}
	return claims, nil
}

func (s *TokenService) DeleteUserTokens(ctx context.Context, userId string) error {
	return s.Repo.DeleteUserTokens(ctx, userId)
}
