package routes

import (
	"github.com/fadhilsurya/mykonsul-mongo/config/config"
	"github.com/fadhilsurya/mykonsul-mongo/internal/handler"
	"github.com/fadhilsurya/mykonsul-mongo/internal/repository"
	"github.com/fadhilsurya/mykonsul-mongo/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializedUser(router *gin.RouterGroup, appConfig config.Config) {
	var (
		userCollection *mongo.Collection   = appConfig.Db.Collection("users")
		userRepo       repository.UserRepo = repository.NewUserRepo(userCollection)

		userService service.UserService = service.NewUserService(userRepo, appConfig)

		userHandler handler.UserHandler = handler.NewUserHandler(userService)
	)

	router.POST("/signup", userHandler.CreateUser)
	router.POST("/signin", userHandler.Login)
	router.DELETE("/", userHandler.DeleteOneUser)

}
