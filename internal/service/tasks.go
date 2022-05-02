package service

import (
	"context"
	"rodnik/internal/apperror"
	"rodnik/internal/entity"
	"rodnik/internal/repository"
	"rodnik/pkg/logger"
)

type taskService struct {
	repTasks repository.Tasks
	repUsers repository.Users
	l        logger.Logger
}

func NewTaskService(repTasks repository.Tasks, repUsers repository.Users, l logger.Logger) *taskService {
	return &taskService{
		repTasks: repTasks,
		repUsers: repUsers,
		l:        l,
	}
}

func (s *taskService) Create(ctx context.Context, task *entity.Task) error {
	user, err := s.repUsers.FindById(ctx, task.CreatorId)
	if err != nil {
		return err
	}
	if user.Leaves < task.Cost {
		//TODO Возможно заранее создать ошибки
		return apperror.PaymentRequired.New(ErrorMessageNoFundsAvailable)
	}

	user.Leaves = user.Leaves - task.Cost
	if err = s.repUsers.UpdateUserBalance(ctx, user); err != nil {
		return err
	}
	if err = s.repTasks.Create(ctx, task); err != nil {
		return err
	}
	return nil
}
