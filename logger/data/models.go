package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

var client *mongo.Client
var db *mongo.Database

func New(mongo *mongo.Client, dbName string) Models {
	client = mongo
	db = client.Database(dbName)

	return Models{LogEntry: LogEntry{}}

}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert() (string, error) {
	collection := db.Collection("logs")

	one, err := collection.InsertOne(context.TODO(), l)
	if err != nil {
		log.Println("failed to insert into logs", err)
		return "", err
	}

	return fmt.Sprint(one.InsertedID), nil

}
