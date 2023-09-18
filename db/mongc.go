package mongc

import (
	"context"
	"fmt"
	"log"
	"time"

	handler "github.com/gin-api/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoClient() *mongo.Client {
	url := handler.GoDotEnvVariable("MONGODB_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}
