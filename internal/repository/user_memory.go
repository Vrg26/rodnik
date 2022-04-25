package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"rodnik/internal/entity"
	"sync"
)

type usersMemoryRepo struct {
	sync.Mutex
	db map[string]*entity.User
}

func NewUsersMemoryRepo() *usersMemoryRepo {
	return &usersMemoryRepo{
		db: make(map[string]*entity.User),
	}
}

func (r usersMemoryRepo) Create(ctx context.Context, user *entity.User) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()
	if r.db == nil {
		r.db = make(map[string]*entity.User)
	}
	user.Id = id.String()
	r.db[user.Id] = user
	return nil
}

func (r usersMemoryRepo) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	for _, value := range r.db {
		if value.Phone == phone {
			return value, nil
		}
	}
	return nil, errors.New("User not found!")
}
