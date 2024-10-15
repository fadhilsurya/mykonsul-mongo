package service

import (
	"context"
	"errors"
	"time"

	"github.com/fadhilsurya/mykonsul-mongo/config/config"
	"github.com/fadhilsurya/mykonsul-mongo/internal/lib/jwt"
	"github.com/fadhilsurya/mykonsul-mongo/internal/model"
	"github.com/fadhilsurya/mykonsul-mongo/internal/repository"
	"github.com/fadhilsurya/mykonsul-mongo/internal/requests"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, req requests.ReqUser) error
	Login(ctx context.Context, email string) (*model.User, *string, error)
	DeleteUser(ctx context.Context, userId string) error
}

type userService struct {
	userRepo repository.UserRepo
	config   config.Config
}

func NewUserService(u repository.UserRepo, conf config.Config) UserService {
	return &userService{
		userRepo: u,
		config:   conf,
	}
}

func (u *userService) CreateUser(ctx context.Context, req requests.ReqUser) error {
	// validation for role thats not user or admin
	if req.Role != "admin" && req.Role != "user" {
		return errors.New("role is invalid")

	}

	data, err := u.userRepo.GetOneUser(ctx, req.Email)
	if err != nil {
		return err
	}

	// validation for double entry
	if data != nil {
		return errors.New("email already registered")
	}

	userModel := model.User{
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		UserId:    uuid.NewString(),
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.userRepo.CreateUser(userModel)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) Login(ctx context.Context, email string) (*model.User, *string, error) {

	data, err := u.userRepo.GetOneUser(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		return nil, nil, nil
	}

	jwt, err := jwt.GenerateJWT(data.UserId, data.Role, u.config.JWTSecret, data.Email)
	if err != nil {
		return nil, nil, err
	}

	return data, &jwt, nil
}

func (u *userService) DeleteUser(ctx context.Context, userId string) error {
	err := u.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}
