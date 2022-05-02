package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"rodnik/internal/apperror"
	"rodnik/internal/entity"
	"sync"
)

type usersMemoryRepo struct {
	sync.Mutex
	db map[uuid.UUID]*entity.User
}

func NewUsersMemoryRepo() *usersMemoryRepo {
	return &usersMemoryRepo{
		db: make(map[uuid.UUID]*entity.User),
	}
}

func (r *usersMemoryRepo) Create(ctx context.Context, user *entity.User) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()
	if r.db == nil {
		r.db = make(map[uuid.UUID]*entity.User)
	}
	user.Id = id
	r.db[user.Id] = user
	return nil
}

func (r *usersMemoryRepo) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	for _, value := range r.db {
		if value.Phone == phone {
			return value, nil
		}
	}
	return nil, apperror.NotFound.New(fmt.Sprintf(ErrorMessageUserNotFoundByPhone, phone))
}
func (r *usersMemoryRepo) FindById(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, ok := r.db[userID]
	if !ok {
		return nil, apperror.NotFound.New(fmt.Sprintf(ErrorMessageUserNotFoundById, userID.String()))
	}
	return user, nil
}
func (r *usersMemoryRepo) UpdateUserBalance(ctx context.Context, user *entity.User) error {
	return nil
}
