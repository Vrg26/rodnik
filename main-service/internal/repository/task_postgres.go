package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"main-service/internal/entity"
	"main-service/pkg/logger"
)

type taskPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

const (
	taskTableName = "tasks"
)

//TODO Информация есть в структуре
var allColumnsTask = []string{"id", "title", "description", "status_id", "creator_id", "helper_id", "cost", "date_relevance", "created_on", "overdue"}

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

func (r *taskPostgres) FindByID(ctx context.Context, taskID uuid.UUID) (*entity.Task, error) {
	task := &entity.Task{}

	sqlText, values, err := squirrel.Select(allColumnsTask...).
		From(taskTableName).
		Where("id = $1", taskID.String()).
		ToSql()

	if err != nil {
		return nil, err
	}
	if err = r.db.GetContext(ctx, task, sqlText, values...); err != nil {
		return nil, err
	}
	return task, nil
}
