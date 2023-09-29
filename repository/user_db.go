package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryDB struct {
	db *mongo.Database
}

func NewUserRepositoryDB(db *mongo.Database) userRepositoryDB {
	return userRepositoryDB{db: db}
}

var colNameUser string = "User"

func (u userRepositoryDB) GetUser(username string) (*User, error) {
	filter := bson.M{"username": username}
	user := User{}
	err := u.db.Collection(colNameUser).FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepositoryDB) CreateUser(user CreateUser) (interface{}, error) {
	collection := u.db.Collection(colNameUser)

	userMap := bson.M{
		"username":  user.UserName,
		"password":  user.Password,
		"email":     user.Email,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}

	res, err := collection.InsertOne(context.Background(), userMap)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}
