package main

import "net/http"

func (app *application) makeRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/healthcheck", makeHTTPHandlerFunc(app.healthcheckHandler))
	router.HandleFunc("POST /v1/movies", makeHTTPHandlerFunc(app.createMovieHandler))
	router.HandleFunc("GET /v1/movies/{id}", makeHTTPHandlerFunc(app.getMovieHandler))
	return router
}
