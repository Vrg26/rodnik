package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"io"
	"main-service/internal/apperror"
	"main-service/internal/entity"
	"main-service/pkg/hash"
	"main-service/pkg/logger"
	"net/http"
	"regexp"
)

//go:generate mockgen -source=users.go -destination=../repository/mocks/user_mock.go -package=mock_repository
type UsersRepo interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	UpdateUserBalance(ctx context.Context, user *entity.User) error
	FindById(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	SetAvatar(ctx context.Context, userID string, avatarName string) error
	AddToFriends(ctx context.Context, friendships *entity.Freindships) error
}

type ClientImageService interface {
	Upload(ctx context.Context, image []byte) (*http.Response, error)
	GetURL() string
}

type UsersService struct {
	client ClientImageService
	repo   UsersRepo
	l      *logger.Logger
}

func NewUserService(client ClientImageService, repo UsersRepo, logger *logger.Logger) *UsersService {
	return &UsersService{
		client: client,
		repo:   repo,
		l:      logger,
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

func (s *UsersService) SetAvatar(ctx context.Context, userID string, imageBytes []byte) (string, error) {
	res, err := s.client.Upload(ctx, imageBytes)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	fileName := string(body)
	if matched, err := regexp.MatchString(".png|.jpg|.jpeg", fileName); err != nil || !matched {
		return "", apperror.Internal.New(fmt.Sprintf(ErrorMessageIncorrectFileName, fileName))
	}

	if err = s.repo.SetAvatar(ctx, userID, fileName); err != nil {
		s.l.Error(err)
		return "", apperror.Internal.New(ErrorMessageInternalServerError)
	}

	return fmt.Sprintf("%s/%s", s.client.GetURL(), string(body)), nil
}

func (s *UsersService) AddToFriends(ctx context.Context, friendships *entity.Freindships) error {

	if friendships.FriendFrom == friendships.FriendTo {
		return apperror.BadRequest.New(ErrorMessageUserCannotAddHimselfAsFriend)
	}

	err := s.repo.AddToFriends(ctx, friendships)
	if err != nil {
		var rqErr *pq.Error
		if errors.As(err, &rqErr) && rqErr.Code == "23503" {
			return apperror.NotFound.New(ErrorMessageFriendNotFound)
		}
		return apperror.Internal.New(ErrorMessageInternalServerError)
	}
	return nil
}
