package v1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"main-service/internal/apperror"
	"main-service/internal/entity"
	"main-service/internal/service"
	"main-service/internal/service/mocks"
	"main-service/pkg/logger"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_taskRoute_create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	type mockTaskServiceBehavior func(s *mock_service.MockTasks, task *entity.Task)
	taskID := uuid.New()
	userID := uuid.New()
	dateRelevance, _ := time.Parse(time.RFC3339, "2022-04-11T07:21:50.32Z")
	createdDate, _ := time.Parse(time.RFC3339, "2022-04-10T07:21:50.32Z")
	testTable := []struct {
		name                    string
		inputBody               string
		inputTask               *entity.Task
		mockTaskServiceBehavior mockTaskServiceBehavior
		expectedStatusCode      int
		expectedResponseBody    string
	}{
		{
			name:      "Ok",
			inputBody: `{"title":"test","description":"test","cost": 200,"date_relevance":"2022-04-11T07:21:50.32Z"}`,
			inputTask: &entity.Task{
				Title:         "test",
				Description:   "test",
				Cost:          200,
				CreatorId:     userID,
				DateRelevance: dateRelevance,
			},
			mockTaskServiceBehavior: func(s *mock_service.MockTasks, task *entity.Task) {
				ctx := context.Background()
				s.EXPECT().Create(ctx, task).Return(&entity.Task{
					Id:            taskID,
					Title:         "test",
					Description:   "test",
					Status:        "OPEN",
					Cost:          200,
					CreatorId:     userID,
					CreatedOn:     createdDate,
					DateRelevance: dateRelevance,
				}, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: fmt.Sprintf(`{"cost":200, "created_on":"2022-04-10T07:21:50.32Z", "creator_id":"%s", "date_relevance":"2022-04-11T07:21:50.32Z", "description":"test", "helper_id":"00000000-0000-0000-0000-000000000000", "id":"%s", "status":"OPEN", "title":"test"}`, userID.String(), taskID.String()),
		},
		{
			name:      "Should return bad request",
			inputBody: `{"description":"test","cost": 200,"date_relevance":"2022-04-11T07:21:50.32Z"}`,
			inputTask: &entity.Task{
				Title:         "test",
				Description:   "test",
				Cost:          200,
				CreatorId:     userID,
				DateRelevance: dateRelevance,
			},
			mockTaskServiceBehavior: func(s *mock_service.MockTasks, task *entity.Task) {},
			expectedStatusCode:      400,
			expectedResponseBody:    `{"message":"field Title is required"}`,
		},
		{
			name:      "Should return payment required",
			inputBody: `{"title":"test","description":"test","cost": 200,"date_relevance":"2022-04-11T07:21:50.32Z"}`,
			inputTask: &entity.Task{
				Title:         "test",
				Description:   "test",
				Cost:          200,
				CreatorId:     userID,
				DateRelevance: dateRelevance,
			},
			mockTaskServiceBehavior: func(s *mock_service.MockTasks, task *entity.Task) {
				ctx := context.Background()
				s.EXPECT().Create(ctx, task).Return(entity.Task{}, apperror.PaymentRequired.New(service.ErrorMessageNoFundsAvailable))
			},
			expectedStatusCode:   402,
			expectedResponseBody: `{"message":"User does not have enough funds to create a task"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			l := logger.New("errors")
			taskService := mock_service.NewMockTasks(c)

			testCase.mockTaskServiceBehavior(taskService, testCase.inputTask)

			taskRoute := &taskRoute{taskService, l}

			//Test Server
			r := gin.Default()
			r.Use(func(c *gin.Context) {
				c.Set("userID", userID.String())
			})
			r.POST("/task", taskRoute.create)

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/task", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
