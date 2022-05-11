package v1

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"main-service/internal/apperror"
	"main-service/internal/entity"
	mock_service "main-service/internal/service/mocks"
	"main-service/pkg/logger"
	"net/http/httptest"
	"testing"
)

func Test_userRoute_addToFriends(t *testing.T) {
	type mockUserServiceBehavior func(s *mock_service.MockUsers, user *entity.Freindships)
	userID := uuid.MustParse("541a283d-043c-447e-b5f2-bcb888a129fc")
	friendID := uuid.MustParse("5879fbfa-ac14-45d7-9e94-9b8ed84a1e77")
	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            *entity.Freindships
		mockUserBehavior     mockUserServiceBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"user_id":"5879fbfa-ac14-45d7-9e94-9b8ed84a1e77"}`,
			inputUser: &entity.Freindships{FriendFrom: userID, FriendTo: friendID},
			mockUserBehavior: func(s *mock_service.MockUsers, user *entity.Freindships) {
				s.EXPECT().AddToFriends(context.Background(), user).Return(nil)
			},
			expectedStatusCode: 201,
		},
		{
			name:                 "Should return bad request",
			inputBody:            `{"user_id":"541a283d"}`,
			inputUser:            &entity.Freindships{FriendFrom: userID, FriendTo: friendID},
			mockUserBehavior:     func(s *mock_service.MockUsers, user *entity.Freindships) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid UUID length: 8"}`,
		},
		{
			name:      "Should return internal error",
			inputBody: `{"user_id":"5879fbfa-ac14-45d7-9e94-9b8ed84a1e77"}`,
			inputUser: &entity.Freindships{FriendFrom: userID, FriendTo: friendID},
			mockUserBehavior: func(s *mock_service.MockUsers, user *entity.Freindships) {
				s.EXPECT().AddToFriends(context.Background(), user).Return(apperror.Internal.New(ErrorMessageInternalServerError))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message": "Internal server error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			logger := logger.New("errors")
			userService := mock_service.NewMockUsers(c)

			testCase.mockUserBehavior(userService, testCase.inputUser)

			userRoute := &userRoute{userService: userService, logger: *logger}

			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("userID", userID.String())
			})
			r.POST("/friends", userRoute.addToFriends)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/friends", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			if w.Code != 201 {
				assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
			}
		})
	}

}
