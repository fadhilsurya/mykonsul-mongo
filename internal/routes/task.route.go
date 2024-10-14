package routes

import (
	"github.com/fadhilsurya/mykonsul-mongo/config/config"
	"github.com/fadhilsurya/mykonsul-mongo/internal/handler"
	"github.com/fadhilsurya/mykonsul-mongo/internal/middleware"
	"github.com/fadhilsurya/mykonsul-mongo/internal/repository"
	"github.com/fadhilsurya/mykonsul-mongo/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializedTask(router *gin.RouterGroup, appConfig config.Config) {
	var (
		taskCollection *mongo.Collection = appConfig.Db.Collection("tasks")
		userCollection *mongo.Collection = appConfig.Db.Collection("users")

		taskRepo repository.TaskRepo = repository.NewTaskRepo(taskCollection)
		userRepo repository.UserRepo = repository.NewUserRepo(userCollection)

		taskService service.TaskService = service.NewTaskService(taskRepo)

		taskHandler handler.TaskHandler = handler.NewTaskHandler(taskService)
	)

	router.POST("/", middleware.MiddlewareToken(appConfig, userRepo), taskHandler.CreateTask)
	router.GET("/:id", middleware.MiddlewareToken(appConfig, userRepo), taskHandler.GetOneTask)
	router.GET("/", middleware.MiddlewareToken(appConfig, userRepo), taskHandler.SearchTask)
	router.PUT("/:id", middleware.MiddlewareToken(appConfig, userRepo), middleware.MiddlewareAdmin(), taskHandler.UpdateTask)
	router.DELETE("/:id", middleware.MiddlewareToken(appConfig, userRepo), middleware.MiddlewareAdmin(), taskHandler.DeleteOneTask)
}
