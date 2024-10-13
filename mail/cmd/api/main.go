package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const webPort = "8080"

func main() {
	app := Config{}

	log.Printf("starting the mail service on port %v\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
