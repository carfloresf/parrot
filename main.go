package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/hellerox/parrot/api"
)

var exit = make(chan os.Signal, 1)

func main() {
	port := os.Getenv("PORT")
	connectionString := os.Getenv("DATABASE_URL")

	if port == "" && connectionString == "" {
		log.Fatalf("missing connection data: %s %s", port, connectionString)
	}

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	a := api.App{}

	a.Initialize(connectionString, port)

	for range exit {
		a.Stop()

		os.Exit(0)
	}
}
