package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connect()

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}

	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Println("Starting front end service on port " + webPort)

	server := http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int = 1
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

		if err != nil {
			counts++
			fmt.Println("Failed to connect to RabbitMQ. Retrying in", backOff)
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Failed to connect to RabbitMQ after 5 retries")
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
	}

	log.Println("Connected to RabbitMQ")
	return connection, nil
}
