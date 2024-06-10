package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}

	defer rabbitConn.Close()

	log.Println("Listening for and consuming messages...")

	consumer, err := event.NewConsumer(rabbitConn, "logs_topic")

	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.ERROR", "log.WARNING"})
	if err != nil {
		log.Fatal("Failed to listen for messages: ", err)
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
