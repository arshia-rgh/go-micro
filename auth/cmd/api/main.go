package main

import (
	"authentication/data"
	"database/sql"
)

const webPort = "8080"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	//app := Config{}

}
