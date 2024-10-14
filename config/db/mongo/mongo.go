package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(mongoAddress string, port string) (*mongo.Database, error) {

	mongoURI := fmt.Sprintf("%s:%s", mongoAddress, port)

	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Replace with your MongoDB URI
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", mongoURI))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal()
		return nil, err
	}

	// Check the connection

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Fatal(err)

	// 	}
	// }()

	connect := client.Database("mykonsul")

	log.Println("Successfully Connect to Mongo DB")

	return connect, nil
}
