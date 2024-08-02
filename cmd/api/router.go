package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) makeRouter() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// movies
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/{id}", app.getMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/{id}", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/{id}", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMovieHandler)
	// user
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPut, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authentication(router)))
}
