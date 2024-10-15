package repository

import (
	"context"
	"errors"

	"github.com/fadhilsurya/mykonsul-mongo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	CreateUser(user model.User) error
	GetOneUser(ctx context.Context, email string) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userRepo struct {
	db *mongo.Collection
}

func NewUserRepo(db *mongo.Collection) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) CreateUser(user model.User) error {
	_, err := u.db.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) GetOneUser(ctx context.Context, email string) (*model.User, error) {
	var (
		user model.User
	)

	err := u.db.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) DeleteUser(ctx context.Context, id string) error {
	_, err := u.db.DeleteOne(ctx, bson.M{"user_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}
