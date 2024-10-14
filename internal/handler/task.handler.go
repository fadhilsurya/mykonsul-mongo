package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/fadhilsurya/mykonsul-mongo/internal/model"
	"github.com/fadhilsurya/mykonsul-mongo/internal/requests"
	"github.com/fadhilsurya/mykonsul-mongo/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	CreateTask(c *gin.Context)
	GetOneTask(c *gin.Context)
	DeleteOneTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	SearchTask(c *gin.Context)
}

type taskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(ts service.TaskService) TaskHandler {
	return &taskHandler{
		taskService: ts,
	}
}

func (ts *taskHandler) CreateTask(c *gin.Context) {
	var (
		req requests.ReqTasks
	)

	ctx := context.Background()

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

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	if u.Role == "user" {
		req.UserId = u.UserId
	}

	err = ts.taskService.CreateTask(ctx, req, req.UserId)
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

func (ts *taskHandler) UpdateTask(c *gin.Context) {
	var (
		req requests.ReqTasksUpdate
	)

	id := c.Param("id")

	if id == "" {
		err := errors.New("id is empty")
		c.AbortWithStatusJSON(http.StatusBadGateway,
			gin.H{"message": "id is empty",
				"status": nil,
				"error":  err.Error()})
		return
	}

	ctx := context.Background()

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	err = ts.taskService.UpdateOneTask(ctx, id, req, req.UserId, req.UserId)
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

func (ts *taskHandler) GetOneTask(c *gin.Context) {

	ctx := context.Background()

	id := c.Param("id")

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

	if u.Role == "user" {
		data, err := ts.taskService.GetOneTask(ctx, id, u.UserId)
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
				"data":   data,
			})
	}

	if u.Role == "admin" {
		data, err := ts.taskService.GetOneTaskAdmin(ctx, id)
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
				"data":   data,
			})
	}
}

func (ts *taskHandler) DeleteOneTask(c *gin.Context) {

	ctx := context.Background()

	id := c.Param("id")

	err := ts.taskService.DeleteOneTask(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"message": "Internal Server Error",
				"status":  nil,
				"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "ok",
			"status": "success",
			"error":  nil,
		})
}

func (ts *taskHandler) SearchTask(c *gin.Context) {

	ctx := context.Background()

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

	perPage := 10 // Default value for perPage
	page := 1     // Default value for page

	// Get perPage from query parameters
	if queryPerPage := c.Query("perPage"); queryPerPage != "" {
		if val, err := strconv.Atoi(queryPerPage); err == nil && val > 0 {
			perPage = val
		}
	}

	// Get page from query parameters
	if queryPage := c.Query("page"); queryPage != "" {
		if val, err := strconv.Atoi(queryPage); err == nil && val > 0 {
			page = val
		}
	}

	skip := (page - 1) * perPage

	// Get title and description from query parameters (optional)
	title := c.Query("title")
	description := c.Query("description")
	userId := c.Query("userId")

	if u.Role == "user" {
		userId = u.UserId
	}

	data, count, err := ts.taskService.GetTasks(ctx, &userId, title, description, perPage, skip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal Server Error",
				"status": nil,
				"error":  err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "ok",
			"status":  "success",
			"error":   nil,
			"data":    data,
			"count":   count,
			"perPage": perPage,
			"page":    skip,
		})

}
