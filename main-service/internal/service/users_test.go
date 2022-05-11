package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"main-service/internal/apperror"
	mock_repository "main-service/internal/repository/mocks"
	mock_image_service "main-service/pkg/client/image_service/mocks"
	"main-service/pkg/logger"
	"net/http"
	"testing"
)

func TestUsersService_SetAvatar(t *testing.T) {
	type mockRepUserBehavior func(s *mock_repository.MockUsersRepo, imageName string)
	type mockClientImageServiceBehavior func(s *mock_image_service.MockClientImageService, response *http.Response)
	userID := uuid.MustParse("5879fbfa-ac14-45d7-9e94-9b8ed84a1e77")
	testTable := []struct {
		name                           string
		fileName                       string
		mockRepUserBehavior            mockRepUserBehavior
		mockClientImageServiceBehavior mockClientImageServiceBehavior
		responseImageService           *http.Response
		expectedError                  error
		expectedURL                    string
	}{
		{
			name:     "ok",
			fileName: "123456.png",
			mockRepUserBehavior: func(s *mock_repository.MockUsersRepo, fileName string) {
				s.EXPECT().SetAvatar(context.Background(), userID.String(), fileName).Return(nil)
			},
			mockClientImageServiceBehavior: func(s *mock_image_service.MockClientImageService, response *http.Response) {
				s.EXPECT().Upload(context.Background(), []byte{}).Return(response, nil)
				s.EXPECT().GetURL().Return("http://localhost")
			},
			responseImageService: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("123456.png")),
			},
			expectedError: nil,
			expectedURL:   "http://localhost/123456.png",
		},
		{
			name:     "Should return error incorrect file name",
			fileName: "123456.png",
			mockRepUserBehavior: func(s *mock_repository.MockUsersRepo, fileName string) {
				s.EXPECT().SetAvatar(context.Background(), userID.String(), fileName).Return(errors.New("some db error"))
			},
			mockClientImageServiceBehavior: func(s *mock_image_service.MockClientImageService, response *http.Response) {
				s.EXPECT().Upload(context.Background(), []byte{}).Return(response, nil)
			},
			responseImageService: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("123456.png")),
			},
			expectedError: apperror.Internal.New(ErrorMessageInternalServerError),
			expectedURL:   "",
		},
		{
			name:                "Should return Internal server error",
			fileName:            "EOF",
			mockRepUserBehavior: func(s *mock_repository.MockUsersRepo, fileName string) {},
			mockClientImageServiceBehavior: func(s *mock_image_service.MockClientImageService, response *http.Response) {
				s.EXPECT().Upload(context.Background(), []byte{}).Return(response, nil)
			},
			responseImageService: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("EOF")),
			},
			expectedError: apperror.Internal.New(fmt.Sprintf(ErrorMessageIncorrectFileName, "EOF")),
			expectedURL:   "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_repository.NewMockUsersRepo(c)
			clientImageService := mock_image_service.NewMockClientImageService(c)
			testCase.mockRepUserBehavior(userRepo, testCase.fileName)
			testCase.mockClientImageServiceBehavior(clientImageService, testCase.responseImageService)

			l := logger.New("errors")

			userService := NewUserService(clientImageService, userRepo, l)

			result, err := userService.SetAvatar(context.Background(), userID.String(), []byte{})

			assert.Equal(t, testCase.expectedURL, result)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
