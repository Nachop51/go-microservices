package main

import (
	"context"
	"log"
	"logger-service/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logger").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Failed to insert log entry:", err)
		return err
	}

	*resp = "Log entry created successfully: " + payload.Name

	return nil
}
