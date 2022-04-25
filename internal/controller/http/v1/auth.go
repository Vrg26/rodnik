package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rodnik/internal/entity"
	"rodnik/internal/service"
	"rodnik/pkg/logger"
)

type loginReq struct {
	Phone    string `json:"phone" binding:"required,min=10"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

type registerReq struct {
	Name     string `json:"name" binding:"required,gte=1,lte=102"`
	Phone    string `json:"phone" binding:"required,min=10"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

type authRoute struct {
	us service.Users
	ts service.Token
	l  *logger.Logger
}

func newAuthRoutes(handler *gin.RouterGroup, us service.Users, ts service.Token, l *logger.Logger) {
	r := &authRoute{us, ts, l}
	h := handler.Group("/auth")
	{
		h.POST("/register", r.register)
		h.POST("/login", r.login)
		h.POST("/logout", AuthUser(ts), r.logout)
		h.POST("/refresh", r.refresh)
	}
}

func (r *authRoute) register(c *gin.Context) {
	ctx := c.Request.Context()
	var regReq registerReq
	if err := c.BindJSON(&regReq); err != nil {
		//todo сделать сообщение об ошибках для людей.
		r.l.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	u := &entity.User{
		Name:     regReq.Name,
		Phone:    regReq.Phone,
		Password: regReq.Password,
	}
	//todo обработать ошибку, если пользователь с таким номером существует
	if err := r.us.Create(ctx, u); err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	tokenPair, err := r.ts.GetTokenPair(ctx, u.Id)
	if err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}
	c.JSON(http.StatusCreated, tokenPair)
}

func (r *authRoute) login(c *gin.Context) {
	ctx := c.Request.Context()

	var lReq loginReq

	if err := c.BindJSON(&lReq); err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	u := &entity.User{
		Phone:    lReq.Phone,
		Password: lReq.Password,
	}
	//todo добавить нормульную обработку ошибок
	if err := r.us.Login(ctx, u); err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Invalid phone and password combination",
		})
		return
	}

	tokenPair, err := r.ts.GetTokenPair(ctx, u.Id)
	if err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}
	c.JSON(http.StatusCreated, tokenPair)
}

func (r *authRoute) logout(c *gin.Context) {
	userId := c.MustGet("userID").(string)
	ctx := c.Request.Context()
	if err := r.ts.DeleteUserTokens(ctx, userId); err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}
	c.Status(http.StatusOK)
}

func (r authRoute) refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var tokenPair *service.TokenPair
	if err := c.BindJSON(&tokenPair); err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	newTokenPair, err := r.ts.RefreshToken(ctx, tokenPair.RefreshToken.String())
	if err != nil {
		r.l.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "refreshToken invalid",
		})
		return
	}
	c.JSON(http.StatusOK, newTokenPair)
}
