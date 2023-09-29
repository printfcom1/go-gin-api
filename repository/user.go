package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName  string             `json:"username"`
	Password  string             `json:"password"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type CreateUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserRepository interface {
	GetUser(string) (*User, error)
	CreateUser(CreateUser) (interface{}, error)
}
