package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"rodnik/internal/apperror"
	"rodnik/internal/entity"
	"rodnik/internal/repository"
	mock_repository "rodnik/internal/repository/mocks"
	"rodnik/pkg/logger"
	"testing"
)

func Test_taskService_Create(t *testing.T) {
	type mockRepTaskBehavior func(s *mock_repository.MockTasks, task *entity.Task)
	type mockRepUserBehavior func(s *mock_repository.MockUsers, userID uuid.UUID)
	testTable := []struct {
		name                string
		inputTask           *entity.Task
		mockRepTaskBehavior mockRepTaskBehavior
		mockRepUserBehavior mockRepUserBehavior
		expectedResult      error
	}{
		{
			name: "Ok",
			inputTask: &entity.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        200,
			},
			mockRepUserBehavior: func(s *mock_repository.MockUsers, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, nil)
				s.EXPECT().UpdateUserBalance(ctx, user).Return(nil)
			},
			mockRepTaskBehavior: func(s *mock_repository.MockTasks, task *entity.Task) {
				ctx := context.Background()
				s.EXPECT().Create(ctx, task).Return(nil)
			},
			expectedResult: nil,
		},
		{
			name: "Should return error \"User has no points\"",
			inputTask: &entity.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        500,
			},
			mockRepUserBehavior: func(s *mock_repository.MockUsers, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, nil)
			},
			mockRepTaskBehavior: func(s *mock_repository.MockTasks, task *entity.Task) {},
			expectedResult:      apperror.PaymentRequired.New(ErrorMessageNoFundsAvailable),
		},
		{
			name: "Repository Failure",
			inputTask: &entity.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        500,
			},
			mockRepUserBehavior: func(s *mock_repository.MockUsers, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, apperror.NotFound.New(repository.ErrorMessageUserNotFoundById))
			},
			mockRepTaskBehavior: func(s *mock_repository.MockTasks, task *entity.Task) {},
			expectedResult:      apperror.NotFound.New(repository.ErrorMessageUserNotFoundById),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			l := logger.New("errors")
			tr := mock_repository.NewMockTasks(c)
			ur := mock_repository.NewMockUsers(c)
			testCase.mockRepTaskBehavior(tr, testCase.inputTask)
			testCase.mockRepUserBehavior(ur, testCase.inputTask.CreatorId)
			ts := NewTaskService(tr, ur, *l)
			ctx := context.Background()

			if result := ts.Create(ctx, testCase.inputTask); result != nil {
				assert.EqualError(t, result, testCase.expectedResult.Error())
			} else {
				assert.Equal(t, result, testCase.expectedResult)
			}
		})
	}
}
