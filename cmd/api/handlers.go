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

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

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

	err = app.models.Movie.Insert(movie)
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

	movie, err := app.models.Movie.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w)
		}
		app.logger.Println(err)
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

	id, err := readParamId(r)
	if err != nil {
		app.CustomErrResponse(w, http.StatusBadGateway, err)
		return
	}

	movie, err := app.models.Movie.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w)
		} else {
			app.serverErrorResponse(w)
		}
		return
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()

	data.ValidateMovie(v, movie)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.Movie.Update(movie)
	if err != nil {
		if errors.Is(err, data.ErrEditConflict) {
			app.editConflictResponse(w)
		} else {
			app.serverErrorResponse(w)
		}
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
	err = app.models.Movie.Delete(id)
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

func (app *application) listMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime",
		"-id", "-title", "-year", "-runtime"}

	data.ValidateFilter(v, input.Filters)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}
	movies, err := app.models.Movie.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w)
		return
	}

	err = WriteJSON(w, http.StatusOK, movies)
	if err != nil {
		app.serverErrorResponse(w)
	}

}
