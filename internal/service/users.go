package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"rodnik/internal/apperror"
	"rodnik/internal/entity"
	"rodnik/pkg/hash"
	"rodnik/pkg/logger"
)

type UsersRepo interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	UpdateUserBalance(ctx context.Context, user *entity.User) error
	FindById(ctx context.Context, userId uuid.UUID) (*entity.User, error)
}

type UsersService struct {
	repo UsersRepo
	l    logger.Logger
}

func NewUserService(repo UsersRepo) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s UsersService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	passHash, err := hash.GetHashAndSalt([]byte(user.Password))
	if err != nil {
		s.l.Error(err)
		return nil, err
	}
	user.Password = passHash

	newUser, err := s.repo.Create(ctx, user)
	if err != nil {
		var rqErr *pq.Error
		if errors.As(err, &rqErr) && rqErr.Code == "23505" {
			return nil, apperror.Conflict.New(ErrorMessageUserIsAlreadyRegistered)
		}
		return nil, err
	}
	return newUser, nil
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
