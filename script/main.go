package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	initTimeZone()
	db := initMongoDB()

	collectionProduct := db.Collection("product")
	keyUniqueProduct := []string{"name", "productCode"}

	createIndexDB(collectionProduct, keyUniqueProduct)

	collectionUser := db.Collection("User")
	keyUniqueUser := []string{"username"}

	createIndexDB(collectionUser, keyUniqueUser)
}

func createIndexDB(collection *mongo.Collection, field []string) {
	for _, key := range field {
		indexModel := mongo.IndexModel{
			Keys:    bson.M{key: 1},
			Options: options.Index().SetUnique(true),
		}

		name, err := collection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			panic(err)
		}
		fmt.Println("create index " + name + "in collection " + collection.Name())
	}
}

func initTimeZone() {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = location
}

func initMongoDB() *mongo.Database {
	url, err := goDotEnvVariable("MONGODB_URL")
	if err != nil {
		panic(err)
	}
	dbName, err := goDotEnvVariable("DB_NAME")
	if err != nil {
		panic(err)
	}
	clientOptions := options.Client().ApplyURI(*url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	mongoDB := client.Database(*dbName)
	fmt.Println("Connected to MongoDB!")
	return mongoDB
}

func goDotEnvVariable(key string) (*string, error) {

	err := godotenv.Load("../.env")

	if err != nil {
		return nil, err
	}

	value := os.Getenv(key)

	return &value, nil
}
