package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rodnik/internal/service"
	"rodnik/pkg/logger"
)

type RConfig struct {
	Logger       *logger.Logger
	UserService  *service.UsersService
	TokenService *service.TokenService
	TaskService  service.Tasks
}

func NewRouter(handler *gin.Engine, c *RConfig) {

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	h := handler.Group("/api")
	{
		newAuthRoutes(h, c.UserService, c.TokenService, c.Logger)
		h.Use(AuthUser(c.TokenService))
		newTaskRoutes(h, c.TaskService, c.Logger)
	}
}
