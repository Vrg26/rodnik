package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main-service/internal/apperror"
	"main-service/internal/entity"
	"main-service/pkg/logger"
	"net/http"
)

type userRoute struct {
	userService UsersService
	logger      logger.Logger
}

type urlResponse struct {
	URL string `json:"url"`
}

type fiendRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

func newUserRoute(handler *gin.RouterGroup, us UsersService, logger logger.Logger) {
	r := userRoute{us, logger}
	h := handler.Group("/users")
	{
		h.POST("/friends", r.addToFriends)
		h.POST("/avatar", r.setAvatar)

	}
}

func (r *userRoute) addToFriends(c *gin.Context) {
	ctx := c.Request.Context()

	userIDFromContext, ok := c.Get("userID")
	if !ok || userIDFromContext == nil {
		sendError(c, apperror.Internal.New(ErrorMessageInternalServerError))
		return
	}
	userID, err := uuid.Parse(userIDFromContext.(string))
	if err != nil {
		sendError(c, apperror.Internal.New(ErrorMessageInternalServerError))
		return
	}

	var fRequest fiendRequest

	if err = c.BindJSON(&fRequest); err != nil {
		sendError(c, err)
		return
	}
	test := fRequest.UserID.String()
	fmt.Println(test)
	friendships := &entity.Freindships{FriendFrom: userID, FriendTo: fRequest.UserID}

	if err = r.userService.AddToFriends(ctx, friendships); err != nil {
		sendError(c, err)
		return
	}
	c.Status(http.StatusCreated)
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
	c.JSON(http.StatusOK, urlResponse{URL: url})
}
