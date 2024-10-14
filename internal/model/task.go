package model

import "time"

type Task struct {
	ID          string    `bson:"_id,omitempty"`
	TaskId      string    `bson:"task_id,omitempty"`
	Title       string    `bson:"title,omitempty"`
	Description string    `bson:"description,omitempty"`
	Status      string    `bson:"status,omitempty"`
	UserId      string    `bson:"user_id,omitempty"`
	CreatedAt   time.Time `bson:"created_at,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty"`
}
