package v1

import (
	"github.com/gin-gonic/gin"
	"main-service/internal/service"
	"main-service/pkg/logger"
	"net/http"
)

type RConfig struct {
	Logger       *logger.Logger
	UserService  UsersService
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
		newUserRoute(h, c.UserService, *c.Logger)
	}
}
