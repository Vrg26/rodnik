package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"rodnik/internal/apperror"
	"rodnik/internal/entity"
	"rodnik/internal/service"
	"rodnik/pkg/logger"
	"time"
)

type taskRoute struct {
	taskService service.Tasks
	l           *logger.Logger
}

type taskRequest struct {
	Title         string    `json:"title,omitempty" binding:"required"`
	Description   string    `json:"description,omitempty" binding:"required"`
	Cost          float64   `json:"cost,omitempty" binding:"required,min=1"`
	DateRelevance time.Time `json:"date_relevance,omitempty" binding:"required" time_format:"2006-01-02T15:04:05.000Z"`
}

func newTaskRoutes(handler *gin.RouterGroup, ts service.Tasks, l *logger.Logger) {
	r := &taskRoute{ts, l}
	h := handler.Group("/task")
	{
		h.POST("/", r.create)
	}
}

func (r taskRoute) create(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == nil {
		returnErrorInResponse(c, apperror.Internal.New(ErrorMessageInternalServerError))
		return
	}

	ctx := c.Request.Context()

	var tReq taskRequest

	if err := c.BindJSON(&tReq); err != nil {
		r.l.Error(err)
		returnErrorInResponse(c, err)
		return
	}
	testTask := &entity.Task{
		Title:         tReq.Title,
		Description:   tReq.Description,
		Cost:          tReq.Cost,
		CreatorId:     userID.(uuid.UUID),
		DateRelevance: tReq.DateRelevance,
	}
	newTask, err := r.taskService.Create(ctx, testTask)

	if err != nil {
		returnErrorInResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, newTask)
}
