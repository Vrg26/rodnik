package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"main-service/internal/apperror"
	"main-service/internal/entity"
	"main-service/pkg/logger"
)

type userPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

const (
	userTableName = "users"
)

func NewUserPostgresRep(db *sqlx.DB, logger logger.Logger) *userPostgres {
	return &userPostgres{db: db, logger: logger}
}

func (r *userPostgres) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	newUser := &entity.User{}

	sqlText, values, err := squirrel.Insert(userTableName).
		Columns("name", "phone", "password", "balance").
		Values(user.Name, user.Phone, user.Password, user.Leaves).
		Suffix("RETURNING \"id\",\"name\", \"phone\", \"password\", \"balance\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	err = r.db.GetContext(ctx, newUser, sqlText, values...)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (r *userPostgres) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	user := &entity.User{}

	sqlText, values, err := squirrel.Select("id", "name", "phone", "password", "balance").
		From(userTableName).
		Where("phone = $1", phone).
		ToSql()

	if err != nil {
		return nil, err
	}

	if err = r.db.GetContext(ctx, user, sqlText, values...); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userPostgres) UpdateUserBalance(ctx context.Context, user *entity.User) error {
	sqlText, values, err := squirrel.Update(userTableName).
		Set("balance", user.Leaves).
		Where("id = $2", user.Id.String()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, sqlText, values...)
	if err != nil {
		return err
	}

	numRow, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRow == 0 {
		return apperror.Internal.New(ErrorMessageBalanceHasNotBeenUpdated)
	}

	return nil
}

func (r *userPostgres) FindById(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	sqlText, values, err := squirrel.Select("id", "name", "phone", "balance").
		From(userTableName).
		Where("id = $1", userID.String()).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err = r.db.GetContext(ctx, user, sqlText, values...); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userPostgres) SetAvatar(ctx context.Context, userID string, avatarName string) error {
	sqlText, values, err := squirrel.Update(userTableName).
		Set("avatar", avatarName).
		Where("id = $2", userID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sqlText, values...)
	if err != nil {
		return err
	}
	return nil
}
