package main

import (
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Println("Starting front end service on port " + webPort)

	server := http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
