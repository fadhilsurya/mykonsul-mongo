package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/fadhilsurya/mykonsul-mongo/internal/model"
	"github.com/fadhilsurya/mykonsul-mongo/internal/requests"
	"github.com/fadhilsurya/mykonsul-mongo/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	Login(c *gin.Context)
	DeleteOneUser(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) UserHandler {
	return &userHandler{
		userService: us,
	}
}

func (us *userHandler) CreateUser(c *gin.Context) {
	var (
		req requests.ReqUser
	)

	ctx := context.Background()

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	err = us.userService.CreateUser(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "ok",
			"status": "success",
			"error":  nil,
		})
}

func (us *userHandler) Login(c *gin.Context) {
	var (
		req requests.ReqLogin
	)

	ctx := context.Background()

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	_, jwt, err := us.userService.Login(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "ok",
			"status": "success",
			"error":  nil,
			"token":  jwt,
		})
}

func (us *userHandler) DeleteOneUser(c *gin.Context) {

	ctx := context.Background()

	// get the user data from JWT
	user, exists := c.Get("user")
	if !exists {
		err := errors.New("user is not exist")
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	u, ok := user.(*model.User)
	if !ok {
		err := errors.New("user is not exist")
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	if u.Role != "admin" {
		err := errors.New("unauthorized user")
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "unauthorized user",
				"status": nil,
				"error":  err.Error()})
		return

	}

	// err := ts.taskService.DeleteOneTask(ctx, id)
	err := us.userService.DeleteUser(ctx, u.UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "ok",
			"status": "success",
			"error":  nil,
		})

}
