package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"main-service/config"
	"main-service/internal/controller/http/v1"
	"main-service/internal/repository"
	"main-service/internal/service"
	"main-service/pkg/client/image-service"
	"main-service/pkg/logger"
	"net/http"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.LogLevel)

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.PG.Host, cfg.PG.Port, cfg.PG.User, cfg.PG.Password, cfg.PG.DBName))

	if err != nil {
		l.Error("Database connection error: %v\n", err)
	}
	clientImageService := image_service.NewClient(cfg.ImageServiceURL, &http.Client{})

	userRepo := repository.NewUserPostgresRep(db, *l)
	tokenRepo := repository.NewTokenMemory()
	taskRepo := repository.NewTaskPostgresRep(db, *l)

	tokenService := service.NewTokenService(tokenRepo, *l, []byte(cfg.SecretKey), 5000, 2000)
	userService := service.NewUserService(clientImageService, userRepo, *l)
	taskService := service.NewTaskService(taskRepo, userRepo, *l)

	rConfig := &v1.RConfig{
		TokenService: tokenService,
		Logger:       l,
		UserService:  userService,
		TaskService:  taskService,
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
