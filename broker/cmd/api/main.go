package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %v\n", webPort)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", webPort),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
