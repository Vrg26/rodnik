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
	h := handler.Group("/tasks")
	{
		h.POST("/", r.create)
	}
}

func (r taskRoute) create(c *gin.Context) {
	userIDctx, ok := c.Get("userID")
	if !ok || userIDctx == nil {
		sendError(c, apperror.Internal.New(ErrorMessageInternalServerError))
		return
	}
	userID, err := uuid.Parse(userIDctx.(string))

	ctx := c.Request.Context()

	var tReq taskRequest

	if err := c.BindJSON(&tReq); err != nil {
		r.l.Error(err)
		sendError(c, err)
		return
	}
	// TODO Работа с статусами
	testTask := &entity.Task{
		Title:         tReq.Title,
		Description:   tReq.Description,
		Cost:          tReq.Cost,
		Status:        "11b812e7-f6c0-4007-b161-b28ca41e5d13",
		CreatorId:     userID,
		DateRelevance: tReq.DateRelevance,
	}
	newTask, err := r.taskService.Create(ctx, testTask)

	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newTask)
}
