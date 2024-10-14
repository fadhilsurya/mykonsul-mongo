package model

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Role      string    `bson:"role"`
	UserId    string    `bson:"user_id"`
	IsActive  bool      `bson:"is_active"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
