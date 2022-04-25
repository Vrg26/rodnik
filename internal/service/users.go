package service

import (
	"context"
	"errors"
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
		return err
	}
	if match := hash.ComparePasswords(u.Password, user.Password); !match {
		return errors.New("Invalid email and password combination")
	}
	*user = *u
	return nil
}
