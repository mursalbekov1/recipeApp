package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {

	var cfg config
	cfg.port = 4000
	cfg.env = "dev"

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, server.Addr)
	err := server.ListenAndServe()
	logger.Fatal(err)
}
