package service

import (
	"context"
	"errors"
	"rodnik/internal/apperror"
	"rodnik/internal/entity"
	"rodnik/internal/repository"
	"rodnik/pkg/hash"
	"rodnik/pkg/logger"
)

type UsersService struct {
	repo repository.Users
	l    logger.Logger
}

func NewUserService(repo repository.Users) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s UsersService) Create(ctx context.Context, user *entity.User) error {
	passHash, err := hash.GetHashAndSalt([]byte(user.Password))
	if err != nil {
		s.l.Error(err)
		return err
	}
	user.Password = passHash
	err = s.repo.Create(ctx, user)
	if err != nil {
		s.l.Error(err)
		return err
	}
	return nil
}

func (s UsersService) Login(ctx context.Context, user *entity.User) error {
	u, err := s.repo.FindByPhone(ctx, user.Phone)
	if err != nil {
		var appError *apperror.AppError
		if errors.As(err, &appError) {
			if appError.Type == apperror.NotFound {
				return apperror.Authorization.New(ErrorMessageInvalidPhoneOrPassword)
			}
		}
		return err
	}
	if match := hash.ComparePasswords(u.Password, user.Password); !match {
		return apperror.Authorization.New(ErrorMessageInvalidPhoneOrPassword)
	}
	*user = *u
	return nil
}
