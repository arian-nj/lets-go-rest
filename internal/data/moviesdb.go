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
	DB *sql.DB
}

// func (s *MovieModel) CreateTable() error {
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	query := `CREATE TABLE IF NOT EXISTS movies (
// 		id bigserial PRIMARY KEY,
// 		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
// 		title text NOT NULL,
// 		year integer NOT NULL,
// 		runtime integer NOT NULL,
// 		genres text[] NOT NULL,
// 		version integer NOT NULL DEFAULT 1
// 		);`

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
// 	_, err = s.db.ExecContext(ctx, query)
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.Commit()
// 	return err
// }

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
	return s.DB.QueryRowContext(ctx, statment, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
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
	err := s.DB.QueryRowContext(ctx, query, id).Scan(args...)
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
	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
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
	result, err := s.DB.ExecContext(ctx, query, id)
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
func (m MovieModel) GetAll(title string, genres []string, filter Filters) ([]*Movie, Metadata, error) {
	var (
		LIMIT  = filter.PageSize
		OFFSET = (filter.Page - 1) * filter.PageSize
	)
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (genres @> $2 OR $2 = '{}')
	ORDER BY %s %s,id ASC LIMIT $3 OFFSET $4`, filter.sortColumn(), filter.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		title,
		pq.Array(genres),
		LIMIT,
		OFFSET,
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totlaRecors := 0
	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&totlaRecors,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totlaRecors, filter.Page, filter.PageSize)
	return movies, metadata, nil
}
