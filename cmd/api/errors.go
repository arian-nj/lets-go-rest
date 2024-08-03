package main

import (
	"errors"
	"net/http"
)

func (app *application) CustomErrResponse(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, envelope{"error": err.Error()})
}

func (app *application) failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, envelope{"error": errors})
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.CustomErrResponse(w, http.StatusNotFound, errors.New("not found"))
}
func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	app.CustomErrResponse(w, http.StatusNotFound, errors.New("method not allowed"))
}
func (app *application) badRequestResponse(w http.ResponseWriter, err error) {
	app.CustomErrResponse(w, http.StatusBadRequest, err)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.CustomErrResponse(w, http.StatusInternalServerError, errors.New("internal server error"))
	app.logger.PrintError(err, nil)
}

func (app *application) editConflictResponse(w http.ResponseWriter) {
	app.CustomErrResponse(w, http.StatusConflict, errors.New("unable to update the record due to an edit conflict, please try again"))
}

func (app *application) toManyRequestsResponse(w http.ResponseWriter) {
	app.CustomErrResponse(w, http.StatusConflict, errors.New("rate limit exceeded"))
}
func (app *application) invalidcredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.CustomErrResponse(w, http.StatusUnauthorized, errors.New(message))
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.CustomErrResponse(w, http.StatusUnauthorized, errors.New(message))
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.CustomErrResponse(w, http.StatusUnauthorized, errors.New(message))
}
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.CustomErrResponse(w, http.StatusForbidden, errors.New(message))

}
