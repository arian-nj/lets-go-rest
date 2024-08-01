package main

import (
	"errors"
	"net/http"

	"github.com/arian-nj/site/back/internal/data"
	"github.com/arian-nj/site/back/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}
	user := data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateUser(v, &user)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "a user with this email already exist")
			app.failedValidationResponse(w, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = WriteJSON(w, http.StatusCreated, user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
