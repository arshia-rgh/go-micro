package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	collection := db.Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		log.Println("failed to get all logs, ", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)

		if err != nil {
			log.Println("error decoding mongodb item to the slice, ", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *LogEntry) GetByID(id string) (*LogEntry, error) {
	collection := db.Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("invalid id format, ", err)
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(context.TODO(), bson.M{"_id": docID}).Decode(&entry)

	if err != nil {
		log.Println("error decoding single log, ", err)
		return nil, err
	}

	return &entry, nil
}
