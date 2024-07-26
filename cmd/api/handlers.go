package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arian-nj/site/back/internal/data"
	"github.com/arian-nj/site/back/internal/validator"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	WriteJSON(w, http.StatusOK, data)
}

type inputMovie struct {
	Title   string       `json:"title"`
	Year    int32        `json:"year"`
	Runtime data.Runtime `json:"runtime"`
	Genres  []string     `json:"genres"`
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input inputMovie
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
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
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.store.InsertMovie(movie)
	if err != nil {
		app.serverErrorResponse(w)
		app.logger.Println(err)
		return
	}
	err = WriteJSON(w, http.StatusOK, input)
	if err != nil {
		app.serverErrorResponse(w)
	}
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readParamId(r)
	if err != nil {
		app.CustomErrResponse(w, http.StatusNotFound, err)
		return
	}

	movie, err := app.store.GetMovieById(id)
	if err != nil {
		app.notFoundResponse(w)
		return
	}

	err = WriteJSON(w, http.StatusOK, envelope{"movie": movie})
	if err != nil {
		app.serverErrorResponse(w)
	}
}

func readParamId(r *http.Request) (int64, error) {
	strId := r.PathValue("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, err
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input inputMovie
	id, err := readParamId(r)
	if err != nil {
		app.CustomErrResponse(w, http.StatusBadGateway, err)
		return
	}

	movie, err := app.store.GetMovieById(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w)
		} else {
			app.serverErrorResponse(w)
		}
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()

	data.ValidateMovie(v, movie)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.store.Update(movie)
	if err != nil {
		app.serverErrorResponse(w)
		app.logger.Println(err)
		return
	}

	err = WriteJSON(w, http.StatusOK, envelope{"movie": movie})
	if err != nil {
		app.serverErrorResponse(w)
	}

}
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readParamId(r)
	if err != nil {
		app.CustomErrResponse(w, http.StatusNotFound, err)
		return
	}
	err = app.store.DeleteMovie(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w)
		} else {
			app.serverErrorResponse(w)
		}
		return
	}

	err = WriteJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"})
	if err != nil {
		app.serverErrorResponse(w)
	}

}
