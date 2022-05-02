package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"rodnik/internal/apperror"
	"rodnik/internal/entity"
	"rodnik/internal/repository"
	"rodnik/internal/service"
	mock_service "rodnik/internal/service/mocks"
	"rodnik/pkg/logger"
	"testing"
)

func Test_authRoute_register(t *testing.T) {
	type mockTokenBehavior func(s *mock_service.MockToken, userId string)
	type mockUserBehavior func(s *mock_service.MockUsers, user *entity.User)
	refreshId, _ := uuid.NewRandom()
	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockTokenBehavior    mockTokenBehavior
		mockUserBehavior     mockUserBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name":"Test","phone":"89229148164","password":"qwert1234"}`,
			inputUser: entity.User{
				Name:     "Test",
				Phone:    "89229148164",
				Password: "qwert1234",
			},
			mockTokenBehavior: func(s *mock_service.MockToken, userID string) {
				ctx := context.Background()
				s.EXPECT().GetTokenPair(ctx, userID).Return(&service.TokenPair{AccessToken: "2345", RefreshToken: refreshId}, nil)
			},
			mockUserBehavior: func(s *mock_service.MockUsers, user *entity.User) {
				ctx := context.Background()
				s.EXPECT().Create(ctx, user).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: fmt.Sprintf(`{"access_token":"2345", "refresh_token":"%s"}`, refreshId.String()),
		},
		{
			name:                 "Empty Fields",
			mockTokenBehavior:    func(s *mock_service.MockToken, userId string) {},
			mockUserBehavior:     func(s *mock_service.MockUsers, user *entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message": "invalid request body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name":"Test","phone":"89229148164","password":"qwert1234"}`,
			inputUser: entity.User{
				Name:     "Test",
				Phone:    "89229148164",
				Password: "qwert1234",
			},
			mockTokenBehavior: func(s *mock_service.MockToken, userId string) {
				ctx := context.Background()
				s.EXPECT().GetTokenPair(ctx, userId).Return(&service.TokenPair{AccessToken: "2345",
					RefreshToken: refreshId}, errors.New("server errors"))
			},
			mockUserBehavior: func(s *mock_service.MockUsers, user *entity.User) {
				ctx := context.Background()
				s.EXPECT().Create(ctx, user).Return(nil)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message": "server errors"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			l := logger.New("errors")
			ts := mock_service.NewMockToken(c)
			us := mock_service.NewMockUsers(c)

			testCase.mockTokenBehavior(ts, testCase.inputUser.Id.String())
			testCase.mockUserBehavior(us, &testCase.inputUser)

			authR := &authRoute{ts: ts, us: us, l: l}

			// Test Server
			r := gin.New()
			r.POST("/register", authR.register)

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register",
				bytes.NewBufferString(testCase.inputBody))

			// Perform Request

			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func Test_authRoute_login(t *testing.T) {
	type mockTokenBehavior func(s *mock_service.MockToken, userId string)
	type mockUserBehavior func(s *mock_service.MockUsers, user *entity.User)
	refreshId, _ := uuid.NewRandom()
	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockTokenBehavior    mockTokenBehavior
		mockUserBehavior     mockUserBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"phone":"89229148164","password":"qwert1234"}`,
			inputUser: entity.User{
				Phone:    "89229148164",
				Password: "qwert1234",
			},
			mockTokenBehavior: func(s *mock_service.MockToken, userID string) {
				ctx := context.Background()
				s.EXPECT().GetTokenPair(ctx, userID).Return(&service.TokenPair{AccessToken: "2345", RefreshToken: refreshId}, nil)
			},
			mockUserBehavior: func(s *mock_service.MockUsers, user *entity.User) {
				ctx := context.Background()
				s.EXPECT().Login(ctx, user).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: fmt.Sprintf(`{"access_token":"2345", "refresh_token":"%s"}`, refreshId.String()),
		},
		{
			name:      "Invalid request body",
			inputBody: `{"password":"qwert1234"}`,
			inputUser: entity.User{
				Phone:    "89229148164",
				Password: "qwert1234",
			},
			mockTokenBehavior:    func(s *mock_service.MockToken, userID string) {},
			mockUserBehavior:     func(s *mock_service.MockUsers, user *entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"field Phone is required"}`,
		},
		{
			name:      "User service Failure",
			inputBody: `{"phone":"89229148164","password":"qwert1234"}`,
			inputUser: entity.User{
				Phone:    "89229148164",
				Password: "qwert1234",
			},
			mockTokenBehavior: func(s *mock_service.MockToken, userID string) {},
			mockUserBehavior: func(s *mock_service.MockUsers, user *entity.User) {
				ctx := context.Background()
				s.EXPECT().Login(ctx, user).Return(apperror.Authorization.New(service.ErrorMessageInvalidPhoneOrPassword))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message": "Invalid phone and password combination"}`,
		},
		{
			name:      "Token service Failure",
			inputBody: `{"phone":"89229148164","password":"qwert1234"}`,
			inputUser: entity.User{
				Phone:    "89229148164",
				Password: "qwert1234",
			},
			mockTokenBehavior: func(s *mock_service.MockToken, userID string) {
				ctx := context.Background()
				s.EXPECT().GetTokenPair(ctx, userID).Return(nil, errors.New("secret key not found"))
			},
			mockUserBehavior: func(s *mock_service.MockUsers, user *entity.User) {
				ctx := context.Background()
				s.EXPECT().Login(ctx, user).Return(nil)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"Internal server error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			ts := mock_service.NewMockToken(c)
			us := mock_service.NewMockUsers(c)
			l := logger.New("errors")

			testCase.mockTokenBehavior(ts, testCase.inputUser.Id.String())
			testCase.mockUserBehavior(us, &testCase.inputUser)

			authR := &authRoute{us, ts, l}

			// Test server
			r := gin.New()
			r.POST("/login", authR.login)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(testCase.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func Test_authRoute_refresh(t *testing.T) {
	type mockTokenBehavior func(s *mock_service.MockToken, refreshToken string)
	refreshId, _ := uuid.NewRandom()
	testTable := []struct {
		name                 string
		inputBody            string
		refreshToken         string
		mockTokenBehavior    mockTokenBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "Ok",
			inputBody:    fmt.Sprintf(`{"access_token":"2345", "refresh_token":"%s"}`, refreshId.String()),
			refreshToken: refreshId.String(),
			mockTokenBehavior: func(s *mock_service.MockToken, refreshToken string) {
				ctx := context.Background()
				s.EXPECT().RefreshToken(ctx, refreshToken).Return(&service.TokenPair{AccessToken: "2345", RefreshToken: refreshId}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"access_token":"2345", "refresh_token":"%s"}`, refreshId.String()),
		},
		{
			name:                 "Empty Fields",
			mockTokenBehavior:    func(s *mock_service.MockToken, refreshToken string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message": "invalid request body"}`,
		},
		{
			name:         "Service Failure",
			refreshToken: refreshId.String(),
			inputBody:    fmt.Sprintf(`{"access_token":"2345", "refresh_token":"%s"}`, refreshId.String()),
			mockTokenBehavior: func(s *mock_service.MockToken, refreshToken string) {
				ctx := context.Background()
				s.EXPECT().RefreshToken(ctx, refreshToken).Return(nil, apperror.Authorization.New(repository.ErrorMessageTokenNotFound))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message": "refreshToken invalid"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			l := logger.New("errors")
			ts := mock_service.NewMockToken(c)

			authR := &authRoute{ts: ts, l: l}

			testCase.mockTokenBehavior(ts, testCase.refreshToken)

			// Test server
			r := gin.New()
			r.POST("/refresh", authR.refresh)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(testCase.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
