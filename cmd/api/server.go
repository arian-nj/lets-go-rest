package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	router := app.makeRouter()

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      router,
		ErrorLog:     log.New(app.logger, "", app.config.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutDownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.logger.PrintInfo("shutting down server ", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutDownErr <- err
		}
		app.logger.PrintInfo("completing background tasks",
			map[string]string{
				"addr": srv.Addr,
			})
		app.wg.Wait()
		shutDownErr <- nil
	}()

	app.logger.PrintInfo(fmt.Sprintf("starting %s server on %s", app.config.env, srv.Addr), nil)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownErr
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
