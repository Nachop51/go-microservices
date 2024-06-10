package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongoDB()

	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register RPC server

	err = rpc.Register(&RPCServer{})

	go app.rpcListen()

	app.serve()
}

func (app *Config) serve() {
	server := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	log.Println("Starting logger service on port", webPort)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server", err)
	}
}

func (app *Config) rpcListen() {
	log.Println("Starting RPC server on port", rpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))

	if err != nil {
		log.Fatal("Failed to start RPC server", err)
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()

		if err != nil {
			log.Println("Failed to accept connection", err)
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongoDB() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error connecting to MongoDB", err)
		return nil, err
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		return nil, err
	}

	return client, nil
}
