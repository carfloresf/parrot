package api

import (
	"fmt"
	"net"

	"github.com/buaazp/fasthttprouter"
	// postgres import
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/hellerox/parrot/controller"
	"github.com/hellerox/parrot/storage"
)

// App represents the application
type App struct {
	listener net.Listener
}

// Initialize sets up the database connection and routes for the app
func (a *App) Initialize(connectionString, appPort string) {
	var controller controller.Controller

	storage := storage.NewStorage(connectionString)
	controller.Service.Storage = storage

	controller.Router = fasthttprouter.New()
	controller.InitializeRoutes()

	server := fasthttp.Server{
		Handler:           fasthttp.CompressHandler(controller.Router.Handler),
		ReadBufferSize:    1024 * 64,
		WriteBufferSize:   1024 * 64,
		ReduceMemoryUsage: true,
	}

	if a.listener != nil {
		log.Fatalf("listener already started")
	}

	var err error

	a.listener, err = net.Listen("tcp4", fmt.Sprint(":", appPort))
	if err != nil {
		log.Fatalf("error creating the listener: %s", err)
	}

	log.Infof("starting server on port %s", appPort)

	go func() {
		err := server.Serve(a.listener)

		if err != nil {
			log.Fatalf("error starting the server: %s", err)
		}
	}()
}

// Stop the API server.
func (a *App) Stop() {
	if a.listener != nil {
		err := a.listener.Close()
		if err != nil {
			log.Fatalf("Error stopping the server: %s", err)
		}

		log.Info("API stopped")
	}
}
