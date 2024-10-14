package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fadhilsurya/mykonsul-mongo/config/config"
	"github.com/fadhilsurya/mykonsul-mongo/internal/lib/jwt"
	"github.com/fadhilsurya/mykonsul-mongo/internal/model"
	"github.com/fadhilsurya/mykonsul-mongo/internal/repository"
	"github.com/gin-gonic/gin"
)

func MiddlewareToken(config config.Config, userRepo repository.UserRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "token is not exist",
					"status": nil,
					"error":  errors.New("token is not exist")})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "bearer token is not exist",
					"status": nil,
					"error":  errors.New("bearer token is not exist")})
			return
		}

		claims, err := jwt.VerifyJWT(tokenString, []byte(config.JWTSecret))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "jwt verify error",
					"status": nil,
					"error":  errors.New("jwt verify error")})
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "wrong format",
					"status": nil,
					"error":  errors.New("wrong format")})
			return
		}

		user, err := userRepo.GetOneUser(ctx, email)
		if err != nil || user == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "user is not found",
					"status": nil,
					"error":  errors.New("user is not found")})
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func MiddlewareAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get("user")
		if !exists {
			err := errors.New("user is not exist")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{"message": "Internal Server Error",
					"status": nil,
					"error":  err.Error()})
			return
		}

		u, ok := user.(*model.User)
		if !ok {
			err := errors.New("user is not exist")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{"message": "Internal Server Error",
					"status": nil,
					"error":  err.Error()})
			return
		}

		if u.Role != "admin" {
			err := errors.New("unauthorized user")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "unauthorized user",
					"status": nil,
					"error":  err.Error()})
			return
		}

		ctx.Next()
	}
}
