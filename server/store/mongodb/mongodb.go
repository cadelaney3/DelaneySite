package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	Client *mongo.Client
	Db     *mongo.Database
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

	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	db := client.Database("delaney-db")

	c := &Client{
		Client: client,
		Db:     db,
	}

	return c, nil
}
