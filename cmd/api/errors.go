package main

import (
	"errors"
	"net/http"
)

func (app *application) CustomErrResponse(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, envelope{"error": err.Error()})
}

func (app *application) failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	WriteJSON(w, http.StatusUnprocessableEntity, envelope{"error": errors})
}

func (app *application) notFoundResponse(w http.ResponseWriter) {
	app.CustomErrResponse(w, http.StatusNotFound, errors.New("not found"))
}
func (app *application) badRequestResponse(w http.ResponseWriter, err error) {
	app.CustomErrResponse(w, http.StatusBadRequest, err)
}

func (app *application) serverErrorResponse(w http.ResponseWriter) {
	app.CustomErrResponse(w, http.StatusInternalServerError, errors.New("internal server error"))
}
