package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"main-service/internal/apperror"
	"main-service/internal/entity"
	"main-service/pkg/hash"
	"main-service/pkg/logger"
)

//go:generate mockgen -source=users.go -destination=../repository/mocks/user_mock.go -package=mock_repository
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

const (
	InitialSum = 1000
)

func (s *UsersService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	passHash, err := hash.GetHashAndSalt([]byte(user.Password))
	if err != nil {
		s.l.Error(err)
		return nil, err
	}
	user.Password = passHash
	user.Leaves = InitialSum
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

func (s *UsersService) Login(ctx context.Context, user *entity.User) error {
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
