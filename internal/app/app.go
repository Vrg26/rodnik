package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rodnik/config"
	v1 "rodnik/internal/controller/http/v1"
	"rodnik/internal/repository"
	"rodnik/internal/service"
	"rodnik/pkg/logger"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.LogLevel)

	userRepo := repository.NewUsersMemoryRepo()
	tokenRepo := repository.NewTokenMemory()

	tokenService := service.NewTokenService(tokenRepo, *l, []byte(cfg.SecretKey), 120, 2000)
	userService := service.NewUserService(userRepo)

	rConfig := &v1.RConfig{
		TokenService: tokenService,
		Logger:       l,
		UserService:  userService,
	}
	handler := gin.Default()
	v1.NewRouter(handler, rConfig)

	//todo запуск сервера вынести
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         ":" + cfg.Port,
	}
	l.Fatal(httpServer.ListenAndServe())
}
