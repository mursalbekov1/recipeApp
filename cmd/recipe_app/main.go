package main

import (
	"github.com/gin-gonic/gin"
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
	router *gin.Engine
}

func main() {

	var cfg config
	cfg.port = 4000
	cfg.env = "dev"

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
		router: gin.Default(),
	}

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.port),
		Handler:      app.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, server.Addr)
	err := server.ListenAndServe()
	logger.Fatal(err)
}
