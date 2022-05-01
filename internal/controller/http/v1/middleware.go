package v1

import (
	"github.com/gin-gonic/gin"
	"rodnik/internal/apperror"
	"rodnik/internal/service"
	"strings"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

func AuthUser(ts service.Token) gin.HandlerFunc {
	return func(c *gin.Context) {
		var h authHeader

		if err := c.ShouldBindHeader(&h); err != nil {
			returnErrorInResponse(c, err)
			return
		}

		tokenHeader := strings.Split(h.Token, "Bearer ")

		if len(tokenHeader) < 2 {
			returnErrorInResponse(c, apperror.Authorization.New(ErrorMessageInvalidHeaderAuth))
			return
		}

		claim, err := ts.ParseToken(tokenHeader[1])
		if err != nil {
			returnErrorInResponse(c, err)
			return
		}
		c.Set("userID", claim.UserID)
		c.Next()
	}
}
