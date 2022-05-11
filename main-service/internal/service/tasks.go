package service

import (
	"context"
	"main-service/internal/apperror"
	"main-service/internal/entity"
	"main-service/pkg/logger"
)

//go:generate mockgen -source=tasks.go -destination=../repository/mocks/task_mock.go -package=mock_repository
type RepoTasks interface {
	Create(ctx context.Context, task *entity.Task) (*entity.Task, error)
}

type taskService struct {
	repTasks RepoTasks
	repUsers UsersRepo
	l        logger.Logger
}

func NewTaskService(repTasks RepoTasks, repUsers UsersRepo, l logger.Logger) *taskService {
	return &taskService{
		repTasks: repTasks,
		repUsers: repUsers,
		l:        l,
	}
}

func (s *taskService) Create(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	user, err := s.repUsers.FindById(ctx, task.CreatorId)
	if err != nil {
		return nil, err
	}
	if user.Leaves < task.Cost {
		//TODO Возможно заранее создать ошибки
		return nil, apperror.PaymentRequired.New(ErrorMessageNoFundsAvailable)
	}

	user.Leaves = user.Leaves - task.Cost
	if err = s.repUsers.UpdateUserBalance(ctx, user); err != nil {
		return nil, err
	}
	newTask, err := s.repTasks.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}
