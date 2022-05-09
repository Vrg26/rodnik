package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"main-service/config"
	"main-service/internal/controller/http/v1"
	repository2 "main-service/internal/repository"
	service2 "main-service/internal/service"
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

	userRepo := repository2.NewUserPostgresRep(db, *l)
	tokenRepo := repository2.NewTokenMemory()
	taskRepo := repository2.NewTaskPostgresRep(db, *l)

	tokenService := service2.NewTokenService(tokenRepo, *l, []byte(cfg.SecretKey), 120, 2000)
	userService := service2.NewUserService(userRepo)
	taskService := service2.NewTaskService(taskRepo, userRepo, *l)

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
