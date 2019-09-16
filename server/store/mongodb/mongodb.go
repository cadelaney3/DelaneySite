package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func InitMongodb(user, password string) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + user + ":" +
		password + "@cluster0-ztlui.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		log.Println("Error with MongoDB client URI: ", err)
		return nil, fmt.Errorf("Error with MongoDB client URI: %s", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Println("Error connecting to MongoDB: ", err)
		return nil, fmt.Errorf("Error connecting to MongoDB: %s", err)
	}

	db := client.Database("delaney-db")

	c := &Client{
		client: client,
		db:     db,
	}

	return c, nil
}
