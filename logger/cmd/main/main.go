package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger/data"
	"os"
	"time"
)

const (
	webPort  = "8080"
	rpcPort  = "5001"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	DB     *mongo.Client
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connectToMongo() (*mongo.Client, error) {

	dbUser := os.Getenv("MONGO_USER")
	dbPassword := os.Getenv("MONGO_PASSWORD")
	dbHost := os.Getenv("MONGO_HOST")
	dbPort := os.Getenv("MONGO_PORT")
	dbName := os.Getenv("MONGO_DB")

	uri := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", dbUser, dbPassword, dbHost, dbPort, dbName)

	clientOptions := options.Client().ApplyURI(uri)

	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB !")

	return mongoClient, nil
}
