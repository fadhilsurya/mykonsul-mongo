package service

import (
	"context"
	"errors"
	"time"

	"github.com/fadhilsurya/mykonsul-mongo/internal/model"
	"github.com/fadhilsurya/mykonsul-mongo/internal/repository"
	"github.com/fadhilsurya/mykonsul-mongo/internal/requests"
	"github.com/google/uuid"
)

type TaskService interface {
	CreateTask(ctx context.Context, req requests.ReqTasks, userId string) error
	GetOneTask(ctx context.Context, id string, userId string) (*model.Task, error)
	UpdateOneTask(ctx context.Context, taskId string, req requests.ReqTasksUpdate, userId string, role string) error
	DeleteOneTask(ctx context.Context, id string) error
	GetOneTaskAdmin(ctx context.Context, id string) (*model.Task, error)
	GetTasks(ctx context.Context, userId *string, title string, description string, perPage int, page int) (*[]model.Task, int, error)
}

type taskService struct {
	taskRepo repository.TaskRepo
}

func NewTaskService(t repository.TaskRepo) TaskService {
	return &taskService{
		taskRepo: t}
}

func (t *taskService) CreateTask(ctx context.Context, req requests.ReqTasks, userId string) error {

	taskModel := model.Task{
		TaskId:      uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Status:      "todo",
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := t.taskRepo.CreateTask(ctx, taskModel)
	if err != nil {
		return err
	}

	return nil
}

func (t *taskService) GetOneTask(ctx context.Context, id string, userId string) (*model.Task, error) {

	data, err := t.taskRepo.GetOneTask(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return data, nil
}

func (t *taskService) GetOneTaskAdmin(ctx context.Context, id string) (*model.Task, error) {

	data, err := t.taskRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return data, nil
}

func (t *taskService) UpdateOneTask(ctx context.Context, taskId string, req requests.ReqTasksUpdate,
	userId string, role string) error {
	var (
		modelTask model.Task
	)

	if req.Status != "done" && req.Status != "in-progress" && req.Status != "todo" {
		return errors.New("request status invalid invalid")
	}

	if req.Title != "" {
		modelTask.Title = req.Title
	}
	if req.Description != "" {
		modelTask.Description = req.Description
	}
	if req.Status != "" {
		modelTask.Status = req.Status
	}

	modelTask.UpdatedAt = time.Now()

	if role == "admin" {
		err := t.taskRepo.UpdateTaskAdmin(ctx, taskId, modelTask)
		if err != nil {
			return err
		}

		return nil
	}

	if role == "user" {
		err := t.taskRepo.UpdateTask(ctx, taskId, modelTask, userId)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("error internal server error")
}

func (t *taskService) DeleteOneTask(ctx context.Context, id string) error {
	err := t.taskRepo.DeleteTask(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *taskService) GetTasks(ctx context.Context, userId *string, title string,
	description string, perPage int, page int) (*[]model.Task, int, error) {

	// Use the repository to get tasks
	tasks, err := t.taskRepo.SearchTasks(ctx, userId, title, description, perPage, page)
	if err != nil {
		return nil, 0, err
	}

	count, err := t.taskRepo.CountTasks(ctx, userId, title, description)
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}
