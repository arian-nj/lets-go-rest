package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/arian-nj/site/back/internal/data"
	"github.com/arian-nj/site/back/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlainText(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.invalidcredentialsResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidcredentialsResponse(w, r)
		return
	}
	token, err := app.models.Token.New(user.ID, 1*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = writeJSON(w, http.StatusCreated, token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
