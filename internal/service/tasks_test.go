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
	type mockRepTaskBehavior func(s *mock_repository.MockRepoTasks, task *entity.Task)
	type mockRepUserBehavior func(s *mock_repository.MockUsersRepo, userID uuid.UUID)
	testTable := []struct {
		name                string
		inputTask           *entity.Task
		mockRepTaskBehavior mockRepTaskBehavior
		mockRepUserBehavior mockRepUserBehavior
		expectedError       error
	}{
		{
			name: "Ok",
			inputTask: &entity.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        200,
			},
			mockRepUserBehavior: func(s *mock_repository.MockUsersRepo, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, nil)
				s.EXPECT().UpdateUserBalance(ctx, user).Return(nil)
			},
			mockRepTaskBehavior: func(s *mock_repository.MockRepoTasks, task *entity.Task) {
				ctx := context.Background()
				s.EXPECT().Create(ctx, task).Return(nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "Should return error \"User has no points\"",
			inputTask: &entity.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        500,
			},
			mockRepUserBehavior: func(s *mock_repository.MockUsersRepo, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, nil)
			},
			mockRepTaskBehavior: func(s *mock_repository.MockRepoTasks, task *entity.Task) {},
			expectedError:       apperror.PaymentRequired.New(ErrorMessageNoFundsAvailable),
		},
		{
			name: "Repository Failure",
			inputTask: &entity.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        500,
			},
			mockRepUserBehavior: func(s *mock_repository.MockUsersRepo, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, apperror.NotFound.New(repository.ErrorMessageUserNotFoundById))
			},
			mockRepTaskBehavior: func(s *mock_repository.MockRepoTasks, task *entity.Task) {},
			expectedError:       apperror.NotFound.New(repository.ErrorMessageUserNotFoundById),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			l := logger.New("errors")
			tr := mock_repository.NewMockRepoTasks(c)
			ur := mock_repository.NewMockUsersRepo(c)
			testCase.mockRepTaskBehavior(tr, testCase.inputTask)
			testCase.mockRepUserBehavior(ur, testCase.inputTask.CreatorId)
			ts := NewTaskService(tr, ur, *l)
			ctx := context.Background()

			_, err := ts.Create(ctx, testCase.inputTask)
			if err != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.Equal(t, err, testCase.expectedError)
			}
		})
	}
}
