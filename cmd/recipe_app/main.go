package main

import (
	"log"
	"net/http"
	"os"
	app "recipeApp/internal/app/api"
	"strconv"
	"time"
)

type config struct {
	port int
	env  string
}

//type Application struct {
//	config config
//	logger *log.Logger
//}

func main() {

	var cfg config
	cfg.port = 4000
	cfg.env = "dev"

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//app := &Application{
	//	config: cfg,
	//	logger: logger,
	//}

	router := app.Routes()

	server := &http.Server{
		Addr:        ":" + strconv.Itoa(cfg.port),
		Handler:     router,
		ReadTimeout: 3 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, server.Addr)
	err := server.ListenAndServe()
	logger.Fatal(err)
}
