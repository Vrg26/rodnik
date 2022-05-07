package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"rodnik/internal/entity"
	"rodnik/pkg/logger"
)

type taskPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

const (
	taskTableName = "tasks"
)

func NewTaskPostgresRep(db *sqlx.DB, logger logger.Logger) *taskPostgres {
	return &taskPostgres{db, logger}
}

func (r *taskPostgres) Create(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	newTask := &entity.Task{}

	sqlText, values, err := squirrel.Insert(taskTableName).
		Columns("title", "description", "status_id", "creator_id", "cost", "date_relevance").
		Values(task.Title, task.Description, task.Status, task.CreatorId, task.Cost, task.DateRelevance).
		Suffix("RETURNING \"id\",\"title\", \"description\", \"status_id\", \"creator_id\", \"cost\", \"date_relevance\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	if err = r.db.GetContext(ctx, newTask, sqlText, values...); err != nil {
		return nil, err
	}
	return newTask, nil
}
