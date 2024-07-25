package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arian-nj/site/back/internal/data"
	"github.com/arian-nj/site/back/internal/validator"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) error {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	WriteJSON(w, http.StatusOK, data)
	return nil
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) error {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	v := validator.New()
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	data.ValidateMovie(v, movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return nil
	}

	err = app.store.InsertMovie(movie)
	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, input)
	return nil
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return WriteJSON(w, http.StatusNotFound, ApiError{Error: err.Error()})

	}

	movie, err := app.store.GetMovieById(id)
	if err != nil {
		return fmt.Errorf("movie by id #%d does not exist", id)
	}
	WriteJSON(w, http.StatusOK, envelope{"movie": movie})
	// fmt.Fprintf(w, "Get Movie #%d", id)
	return nil
}

func GetId(r *http.Request) (int64, error) {
	strId := r.PathValue("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, err
}
