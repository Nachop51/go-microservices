package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection, queueName string) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()

	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	ch, err := c.conn.Channel()

	if err != nil {
		return err
	}

	return declareExchange(ch)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareRandomQueue(ch)

	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(
			q.Name,
			topic,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload

			err := json.Unmarshal(d.Body, &payload)

			if err != nil {
				fmt.Println("Failed to unmarshal message")
			}

			fmt.Printf("Received a message: %s\n", payload.Data)
			go handlePayload(payload)
		}
	}()

	fmt.Println("Waiting for messages [Exchange, Queue] [logs_topic, " + q.Name + "]")
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)

		if err != nil {
			fmt.Println("Failed to log event")
		}

		fmt.Println("Event logged")
	default:
		fmt.Println("Unknown message type")
	}
}

func logEvent(data Payload) error {
	jsonData, _ := json.MarshalIndent(data, "", "\t")

	response, err := http.Post("http://logger-service/log", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
