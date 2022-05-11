package v1

import (
	"github.com/gin-gonic/gin"
	"main-service/internal/apperror"
	"main-service/internal/service"
	"strings"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

func AuthUser(ts service.Token) gin.HandlerFunc {
	return func(c *gin.Context) {
		var h authHeader

		if err := c.ShouldBindHeader(&h); err != nil {
			sendError(c, err)
			return
		}

		tokenHeader := strings.Split(h.Token, "Bearer ")

		if len(tokenHeader) < 2 {
			sendError(c, apperror.Authorization.New(ErrorMessageInvalidHeaderAuth))
			return
		}

		claim, err := ts.ParseToken(tokenHeader[1])
		if err != nil {
			sendError(c, err)
			return
		}
		c.Set("userID", claim.UserID)
		c.Next()
	}
}
