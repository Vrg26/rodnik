package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.PG.Host, cfg.PG.Port, cfg.PG.User, cfg.PG.Password, cfg.PG.DBName))

	if err != nil {
		l.Error("Database connection error: %v\n", err)
	}

	userRepo := repository.NewUserPostgresRep(db, *l)
	tokenRepo := repository.NewTokenMemory()
	taskRepo := repository.NewTaskPostgresRep(db, *l)

	tokenService := service.NewTokenService(tokenRepo, *l, []byte(cfg.SecretKey), 120, 2000)
	userService := service.NewUserService(userRepo)
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
