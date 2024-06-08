package main

import (
	"authentication/data"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "80"

type Config struct {
	db     *pgxpool.Pool
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	con := connectToDB()

	if con == nil {
		log.Fatal("Error connecting to database")
	}

	app := Config{
		db:     con,
		Models: data.New(con),
	}

	server := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server", err)
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	log.Println("dsn", dsn)
	db, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())

	log.Println("Ping")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Connected to database")

	return db, nil
}

func connectToDB() *pgxpool.Pool {
	dsn := os.Getenv("DSN")
	var counts int = 0

	for {
		db, err := openDB(dsn)

		if err != nil {
			log.Println("Error connecting to database, may be not ready yet")
			counts++
		} else {
			return db
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Retrying in 3 seconds")
		time.Sleep(3 * time.Second)
	}
}
