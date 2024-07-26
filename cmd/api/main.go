package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arian-nj/site/back/internal/data"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	models *data.Models
}

const version = "1.0.0"

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.Parse()

	app := application{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		config: cfg,
	}

	models, err := data.NewModels()
	if err != nil {
		app.logger.Panic(err)
	}
	app.models = models

	router := app.makeRouter()

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("starting %s server on %s", app.config.env, srv.Addr)
	err = srv.ListenAndServe()
	app.logger.Fatal(err)

}
