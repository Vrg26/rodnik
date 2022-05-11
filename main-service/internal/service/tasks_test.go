package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"main-service/internal/apperror"
	entity2 "main-service/internal/entity"
	"main-service/internal/repository"
	mock_repository2 "main-service/internal/repository/mocks"
	"main-service/pkg/logger"
	"testing"
)

func Test_taskService_Create(t *testing.T) {
	type mockRepTaskBehavior func(s *mock_repository2.MockRepoTasks, task *entity2.Task)
	type mockRepUserBehavior func(s *mock_repository2.MockUsersRepo, userID uuid.UUID)
	testTable := []struct {
		name                string
		inputTask           *entity2.Task
		mockRepTaskBehavior mockRepTaskBehavior
		mockRepUserBehavior mockRepUserBehavior
		expectedError       error
	}{
		{
			name: "Ok",
			inputTask: &entity2.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        200,
			},
			mockRepUserBehavior: func(s *mock_repository2.MockUsersRepo, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity2.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, nil)
				s.EXPECT().UpdateUserBalance(ctx, user).Return(nil)
			},
			mockRepTaskBehavior: func(s *mock_repository2.MockRepoTasks, task *entity2.Task) {
				ctx := context.Background()
				s.EXPECT().Create(ctx, task).Return(nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "Should return error \"User has no points\"",
			inputTask: &entity2.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        500,
			},
			mockRepUserBehavior: func(s *mock_repository2.MockUsersRepo, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity2.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, nil)
			},
			mockRepTaskBehavior: func(s *mock_repository2.MockRepoTasks, task *entity2.Task) {},
			expectedError:       apperror.PaymentRequired.New(ErrorMessageNoFundsAvailable),
		},
		{
			name: "Repository Failure",
			inputTask: &entity2.Task{
				Title:       "test",
				Description: "test",
				CreatorId:   uuid.New(),
				Cost:        500,
			},
			mockRepUserBehavior: func(s *mock_repository2.MockUsersRepo, userID uuid.UUID) {
				ctx := context.Background()
				user := &entity2.User{Leaves: 200, Name: "test", Id: userID}
				s.EXPECT().FindById(ctx, userID).Return(user, apperror.NotFound.New(repository.ErrorMessageUserNotFoundById))
			},
			mockRepTaskBehavior: func(s *mock_repository2.MockRepoTasks, task *entity2.Task) {},
			expectedError:       apperror.NotFound.New(repository.ErrorMessageUserNotFoundById),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			l := logger.New("errors")
			tr := mock_repository2.NewMockRepoTasks(c)
			ur := mock_repository2.NewMockUsersRepo(c)
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

func Test_taskService_GetByID(t *testing.T) {
	type mockRepTaskBehavior func(s *mock_repository2.MockRepoTasks, task *entity2.Task)
	taskID := uuid.New()
	testTable := []struct {
		name                string
		inputID             string
		mockRepTaskBehavior mockRepTaskBehavior
		expectedValue       *entity2.Task
		expectedError       error
	}{
		{
			name:    "OK",
			inputID: taskID.String(),
			expectedValue: &entity2.Task{
				Id:          taskID,
				Title:       "test",
				Description: "test",
			},
			//mockRepTaskBehavior: func(s *mock_repository.MockRepoTasks, task *entity.Task) {
			//	s.EXPECT().GetByID(context.Background(), taskID)
			//},
			expectedError: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

		})
	}
}
