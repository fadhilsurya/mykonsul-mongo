package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadhilsurya/mykonsul-mongo/config/config"
	"github.com/fadhilsurya/mykonsul-mongo/config/db/redis"
	"github.com/fadhilsurya/mykonsul-mongo/internal/middleware"
	"github.com/fadhilsurya/mykonsul-mongo/internal/routes"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {

	// db.MigrationData()

	config.InitConfig()
	serverConfig := config.AppConfig

	redis := redis.RedisClient(serverConfig.RedisConfig.Address, serverConfig.RedisConfig.Port)

	ginMode := serverConfig.GinMode
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	r := gin.Default()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.App.Port),
		Handler: r,
	}

	// set rame limit per 1 msg per 5 burst
	limiter := rate.NewLimiter(1, 5)
	r.Use(middleware.RateLimitMiddleware(limiter))

	routes.InitializeRoutes(r, &serverConfig, redis)

	go func() {
		fmt.Printf("listening to port : %v", serverConfig.App.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe: %v", err)
		}
	}()

	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println("----shutdown server----")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}

	log.Println("Server exiting")
}
