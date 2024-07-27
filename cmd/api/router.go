package main

import "net/http"

func (app *application) makeRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	router.HandleFunc("POST /v1/movies", app.createMovieHandler)
	router.HandleFunc("GET /v1/movies/{id}", app.getMovieHandler)
	router.HandleFunc("PATCH /v1/movies/{id}", app.updateMovieHandler)
	router.HandleFunc("DELETE /v1/movies/{id}", app.deleteMovieHandler)
	router.HandleFunc("GET /v1/movies", app.listMovieHandler)
	return router
}
