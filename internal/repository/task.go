package repository

import (
	"context"
	"errors"

	"github.com/fadhilsurya/mykonsul-mongo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, task model.Task) error
	GetOneTask(ctx context.Context, id string, userId string) (*model.Task, error)
	GetOne(ctx context.Context, taskId string) (*model.Task, error)
	UpdateTask(ctx context.Context, id string, task model.Task, userId string) error
	UpdateTaskAdmin(ctx context.Context, id string, task model.Task) error
	DeleteTask(ctx context.Context, id string) error
	CountTasks(ctx context.Context, userId *string, title string, description string) (int, error)
	SearchTasks(ctx context.Context, userId *string, title string,
		description string, perPage int, page int) (*[]model.Task, error)
}

type taskRepo struct {
	db *mongo.Collection
}

func NewTaskRepo(db *mongo.Collection) TaskRepo {
	return &taskRepo{
		db: db,
	}
}

func (t *taskRepo) CreateTask(ctx context.Context, task model.Task) error {
	_, err := t.db.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepo) SearchTasks(ctx context.Context, userId *string, title string,
	description string, perPage int, page int) (*[]model.Task, error) {
	var tasks []model.Task

	filter := bson.M{}

	if userId != nil {
		filter["user_id"] = *userId
	}

	if title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"} // Case-insensitive search in title
	}

	if description != "" {
		filter["description"] = bson.M{"$regex": description, "$options": "i"} // Case-insensitive search in description
	}

	// Find tasks with pagination
	cursor, err := t.db.Find(ctx, filter, options.Find().SetLimit(int64(perPage)).SetSkip(int64(page)))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	// Decode the cursor into a slice of tasks
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return &tasks, nil
}

func (t *taskRepo) CountTasks(ctx context.Context, userId *string, title string, description string) (int, error) {
	filter := bson.M{}

	if userId != nil {
		filter["user_id"] = *userId
	}

	if title != "" {
		// Case-insensitive search in title
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}

	if description != "" {
		// Case-insensitive search in description
		filter["description"] = bson.M{"$regex": description, "$options": "i"}
	}

	totalCount, err := t.db.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(totalCount), nil
}

func (t *taskRepo) GetOneTask(ctx context.Context, id string, userId string) (*model.Task, error) {
	var (
		task model.Task
	)

	err := t.db.FindOne(ctx, bson.M{
		"task_id": id,
		"user_id": userId,
	}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (t *taskRepo) GetOne(ctx context.Context, taskId string) (*model.Task, error) {
	var (
		task model.Task
	)

	err := t.db.FindOne(ctx, bson.M{
		"task_id": taskId,
	}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (t *taskRepo) UpdateTask(ctx context.Context, id string, task model.Task, userId string) error {
	_, err := t.db.UpdateOne(ctx, bson.M{"task_id": id, "user_id": userId}, bson.M{"$set": task})
	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepo) UpdateTaskAdmin(ctx context.Context, id string, task model.Task) error {

	_, err := t.db.UpdateOne(ctx, bson.M{"task_id": id}, bson.M{"$set": task})
	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepo) DeleteTask(ctx context.Context, id string) error {
	_, err := t.db.DeleteOne(ctx, bson.M{"task_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("task not found")
		}
		return err
	}

	return nil
}
