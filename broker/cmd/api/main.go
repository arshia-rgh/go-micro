package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitMQ, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	defer rabbitMQ.Close()

	app := Config{
		Rabbit: rabbitMQ,
	}

	log.Printf("Starting broker service on port %v\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var connection *amqp.Connection
	var counts int64
	var backOff = 1 * time.Second

	rabbitMQURL := fmt.Sprintf("amqp://%v:%v@%v:%v/",
		os.Getenv("RABBITMQ_USERNAME"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	for {
		c, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Println("rabbitmq not yet ready...!")
			counts++
		} else {
			log.Println("connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("backing off for %v seconds\n", backOff)
		time.Sleep(backOff)

	}

	return connection, nil
}
