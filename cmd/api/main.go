package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{err.Error()})
		}
	}
}

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	store  storage
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

	pgstore, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer pgstore.db.Close(context.Background())

	app.store = pgstore
	err = app.store.Init()
	if err != nil {
		log.Fatal(err)
	}

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
