package data

import (
	"errors"
	"time"

	"github.com/arian-nj/site/back/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"-"  db:"created_at"`
	Title     string    `json:"title" db:"title"`
	Year      int32     `json:"year"  db:"year"`
	Runtime   Runtime   `json:"runtime"  db:"runtime"`
	Genres    []string  `json:"genres"  db:"genres"`
	Version   int32     `json:"version"  db:"version"`
}

var (
	ErrRecordNotFound = errors.New("record not found")
)

// func NewMovie(title string, year int32, runtime Runtime, genres []string) *Movie {
// 	return &Movie{
// 		Title:     title,
// 		Year:      year,
// 		Runtime:   runtime,
// 		Genres:    genres,
// 		CreatedAt: time.Now().UTC(),
// 	}
// }

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

}
