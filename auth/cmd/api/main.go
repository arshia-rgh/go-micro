package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "8080"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load the env file: %v", err)
		return
	}
	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	dbName := os.Getenv("dbName")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panicf("Could not connect to the DB: %v", err)
	}
	models := data.New(db)

	app := Config{
		DB:     db,
		Models: models,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load the env file: %v", err)
		return nil
	}
	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	dbName := os.Getenv("dbName")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres !")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Try again in two seconds ...")
		time.Sleep(2 * time.Second)

	}
}
