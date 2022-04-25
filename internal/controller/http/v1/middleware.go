package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		tokenHeader := strings.Split(h.Token, "Bearer ")

		if len(tokenHeader) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Must provide Authorization header with format 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		claim, err := ts.ParseToken(tokenHeader[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "failed to parse token",
			})
			c.Abort()
			return
		}
		c.Set("userID", claim.UserID)
		c.Next()
	}
}
