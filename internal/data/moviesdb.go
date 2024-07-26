package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type MovieModel struct {
	db *sql.DB
}

func (s *MovieModel) CreateTable() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `CREATE TABLE IF NOT EXISTS movies (
		id bigserial PRIMARY KEY,
		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		title text NOT NULL,
		year integer NOT NULL,
		runtime integer NOT NULL,
		genres text[] NOT NULL,
		version integer NOT NULL DEFAULT 1
		);`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (s *MovieModel) Insert(movie *Movie) error {
	fmt.Println("making ", movie.Title, " in db")
	statment := `INSERT INTO movies 
	(title,year,runtime,genres)
	VALUES 
	($1,$2,$3,$4) 
	RETURNING id, created_at, version
	`
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return s.db.QueryRowContext(ctx, statment, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (s *MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`
	var movie Movie

	args := []interface{}{
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, id).Scan(args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &movie, nil
}

func (s *MovieModel) Update(movie *Movie) error {
	query := `
	UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version =
	version + 1
	WHERE id = $5 AND version=$6
	RETURNING version`
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (s *MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM movies
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsEffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
