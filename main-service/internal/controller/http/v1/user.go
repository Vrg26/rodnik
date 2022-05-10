package v1

import (
	"github.com/gin-gonic/gin"
	"main-service/internal/apperror"
	"main-service/pkg/logger"
	"net/http"
)

type userRoute struct {
	userService UsersService
	logger      logger.Logger
}

func newUserRoute(handler *gin.RouterGroup, us UsersService, logger logger.Logger) {
	r := userRoute{us, logger}
	h := handler.Group("/users")
	{
		h.POST("/avatar", r.setAvatar)
	}
}

func (r *userRoute) setAvatar(c *gin.Context) {
	userId := c.MustGet("userID").(string)
	ctx := c.Request.Context()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5<<20)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		sendError(c, apperror.BadRequest.New(err.Error()))
		return
	}
	defer file.Close()
	buffer := make([]byte, fileHeader.Size)

	file.Read(buffer)
	fileType := http.DetectContentType(buffer)
	if fileType != "image/jpeg" && fileType != "image/png" {
		sendError(c, apperror.BadRequest.New("file type is not supported"))
		return
	}
	url, err := r.userService.SetAvatar(ctx, userId, buffer)
	c.JSON(http.StatusOK, url)
}
