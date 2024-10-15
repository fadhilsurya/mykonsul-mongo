package routes

import (
	"github.com/fadhilsurya/mykonsul-mongo/config/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitializeRoutes(router *gin.Engine, appConfig *config.Config, rds *redis.Client) {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PONGG!",
			"error":   nil,
			"data":    nil,
		})
	})

	taskGroup := router.Group("/tasks")
	userGroup := router.Group("/users")

	InitializedTask(taskGroup, *appConfig, rds)
	InitializedUser(userGroup, *appConfig)

}
