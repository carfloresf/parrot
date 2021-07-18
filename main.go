package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/hellerox/parrot/api"
)

var exit = make(chan os.Signal, 1) // nolint: gochecknoglobals

func main() {
	a := api.App{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		connectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			"bird", "docker", "localhost", 5433, "parrot")
	}

	log.Infof("setting connectionString: %s", connectionString)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	err := a.Initialize(connectionString, port)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for range exit {
		log.Info("Stopping server...")

		err := a.Stop()
		if err != nil {
			log.Fatalf("Error stopping the server: %s", err)
		}

		os.Exit(0)
	}
}
