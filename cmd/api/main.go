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
	"github.com/arian-nj/site/back/internal/jsonlog"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type config struct {
	port    int
	env     string
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models *data.Models
}

const version = "1.0.0"

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	l := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	app := application{
		logger: l,
		config: cfg,
	}

	models, err := data.NewModels()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
	app.models = models

	router := app.makeRouter()

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      router,
		ErrorLog:     log.New(app.logger, "", cfg.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.PrintInfo(fmt.Sprintf("starting %s server on %s", app.config.env, srv.Addr), nil)

	err = srv.ListenAndServe()
	app.logger.PrintFatal(err, nil)

}
