package main

import "log"

type Config struct {
}

const webPort = "8080"

func main() {
	app := Config{}

	log.Printf("starting the mail service on port %v\n", webPort)
}
