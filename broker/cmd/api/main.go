package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %v\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
