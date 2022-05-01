package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"rodnik/internal/apperror"
	"rodnik/internal/service"
	mock_service "rodnik/internal/service/mocks"
	"testing"
)

func TestAuthUser(t *testing.T) {
	type mockTokenBehavior func(s *mock_service.MockToken, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockTokenBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockToken, token string) {
				s.EXPECT().ParseToken(token).Return(&service.CustomClaims{UserID: "1"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "No Header",
			headerName:           "",
			mockBehavior:         func(s *mock_service.MockToken, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message": "Must provide Authorization header with format 'Bearer {token}'"}`,
		},
		{
			name:                 "Invalid Bearer",
			headerName:           "Authorization",
			headerValue:          "Beer token",
			token:                "token",
			mockBehavior:         func(s *mock_service.MockToken, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message": "Must provide Authorization header with format 'Bearer {token}'"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockToken, token string) {
				s.EXPECT().ParseToken(token).Return(&service.CustomClaims{UserID: "1"},
					apperror.Authorization.New(service.ErrorMessageInvalidToken))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"Invalid token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			ts := mock_service.NewMockToken(c)
			testCase.mockBehavior(ts, testCase.token)

			// Test Server
			r := gin.New()
			r.GET("/protected", AuthUser(ts), func(c *gin.Context) {
				id, _ := c.Get("userID")
				c.String(200, id.(string))
			})

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			//Make request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
